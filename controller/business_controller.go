package controller

import (
	"fmt"
	"log"
	"net/http"
	"smart-queue/model"
	"smart-queue/sse"
	"smart-queue/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessController struct {
	Broadcaster *sse.Broadcaster
	DB          *gorm.DB
}

func (bc *BusinessController) ListBusiness(c *gin.Context) {
	var businesses []model.Business
	category := c.Query("category")

	query := bc.DB.Model(&model.Business{})

	if category != "" {
		query = query.Where("type= ?", category).Select("id, name, type, rating")
	}

	log.Println(category)
	if err := query.Find(&businesses).Error; err != nil {
		utils.ResponseError(c, 500, "failed to list businesses")
		return
	}

	utils.ResponseSuccess(c, "Fetched Successfully", businesses)
}

func (bc *BusinessController) FetchBusinessDetailById(c *gin.Context) {

	var businessDetail model.Business

	var id = c.Param("id")

	log.Println("id", id)

	if err := bc.DB.First(&businessDetail, id).Error; err != nil {
		utils.ResponseError(c, 400, "failed to fetch Data")
		return
	}

	utils.ResponseSuccess(c, "Data feched succesfully", businessDetail)

}

func (h *BusinessController) StreamSSE(c *gin.Context) {
	businessIDParam := c.Param("businessID")
	var businessID uint
	fmt.Sscanf(businessIDParam, "%d", &businessID)

	ch := h.Broadcaster.Subscribe(businessID)
	defer h.Broadcaster.Unsubscribe(businessID, ch)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)

	notify := c.Writer.CloseNotify()

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	// for ev := range ch {
	// 	fmt.Fprintf(c.Writer, "data: %s\n\n", ev.JSON())
	// 	c.Writer.Flush()
	// }

	for {
		select {
		case <-notify:
			return

		case ev := <-ch:

			// return
			fmt.Fprintf(c.Writer, "data: %s\n\n", ev.JSON())
		}

		flusher.Flush()
	}
}

// func (h *BusinessController) BroadcastMessage(c *gin.Context) {
// 	var req struct {
// 		Message string `json:"message"`
// 	}
// 	c.BindJSON(&req)

// 	businessIDParam := c.Param("businessID")
// 	var businessID uint
// 	fmt.Sscanf(businessIDParam, "%d", &businessID)

// 	h.Broadcaster.Publish(businessID, sse.Event{
// 		Type: "custom_message",
// 		Data: req.Message,
// 	})

// 	c.JSON(http.StatusOK, gin.H{"status": "sent"})
// }

func (h *BusinessController) JoinQueue(c *gin.Context) {
	bid, _ := strconv.ParseInt(c.Param("business_id"), 10, 64)

	var req struct {
		CustomerId int64 `json:"customer_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	var existing model.Queue
	err := h.DB.Where("business_id = ? AND customer_id = ? AND status IN ?",
		bid, req.CustomerId, []string{"waiting", "called"}).
		First(&existing).Error

	if err == nil {
	
		c.JSON(http.StatusConflict, gin.H{
			"error":  "Customer already in queue for this business",
			"queue":  existing,
			"status": false,
		})
		return
	}

	if err != gorm.ErrRecordNotFound {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}


	var position int64
	if err := h.DB.Model(&model.Queue{}).
		Where("business_id = ? AND status IN ?", bid, []string{"waiting", "called"}).
		Count(&position).Error; err != nil {
		log.Println("Query Error:", err)
	}

	number := position + 1
	now := time.Now()

	const avgServiceTime = 5
	estimatedWait := number * avgServiceTime
	wait := int(estimatedWait)
	queue := model.Queue{
		BusinessID:        bid,
		CustomerID:        req.CustomerId,
		QueueNumber:       fmt.Sprintf("Q%d", number),
		Status:            "waiting",
		PartySize:         1,
		Position:          number,
		JoinedAt:          &now,
		EstimatedWaitTime: (&wait),
	}

	if err := h.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join queue"})
		return
	}

	h.Broadcaster.Subscribe(uint(bid))
	c.JSON(http.StatusCreated, queue)
}

func (bc *BusinessController) LeaveQueue(c *gin.Context) {
	businessId, err := strconv.Atoi(c.Param("business_id"))

	if err != nil || businessId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business_id"})
		return
	}

	var req struct {
		CustomerId int64 `json:"customer_id"`
	}


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := bc.DB.Begin()

	if tx.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could Not Do the Transaction"})
		return
	}

	var entry model.Queue
	if err := tx.Where(
		"business_id = ? AND customer_id = ? AND status = ?",
		businessId, req.CustomerId, "waiting",
	).First(&entry).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Queue entry not found"})
		return
	}

	now := time.Now()
	entry.Status = "left"
	entry.CompletedAt = &now

	if err := tx.Save(&entry).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update queue entry"})
		return
	}


	var avgServiceTime float64
	if err := tx.Raw(`
    SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - joined_at))/60), 0)
    FROM queues
    WHERE business_id = ? AND status = 'completed' AND completed_at IS NOT NULL
`, businessId).Scan(&avgServiceTime).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate average time", "details": err.Error()})
		return
	}

	if avgServiceTime <= 0 {
		avgServiceTime = 5 
	}


	if err := tx.Exec(`
        UPDATE queues
        SET 
            position = position - 1,
            estimated_wait_time = CAST(((position - 1) * ?) AS INTEGER)
        WHERE business_id = ?
          AND position > ?
          AND status IN ('waiting', 'called')
    `, avgServiceTime, businessId, entry.Position).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shift positions"})
		return
	}


	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit failed"})
		return
	}


	bc.Broadcaster.Subscribe(uint(businessId))

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Queue left successfully",
	})

}

//CallNext

func (bc *BusinessController) CallNext(c *gin.Context) {

	businessId, err := strconv.Atoi(c.Param("business_id"))

	if err != nil || businessId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business_id"})
		return
	}

	var req struct {
		CustomerId int64 `json:"customer_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer not found"})
		return
	}

	tx := bc.DB.Begin()

	var entry model.Queue
	if err := tx.Where("business_id = ? AND customer_id = ? AND status = ?", businessId, req.CustomerId, "waiting").First(&entry).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Queue entry not found"})
		return
	}

	now := time.Now()
	entry.Status = "called"
	entry.CompletedAt = &now
	entry.Position = 0

	if err := tx.Save(&entry).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update queue entry"})
		return
	}

	var avgServiceTime float64
	if err := tx.Raw(`
    SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - joined_at))/60), 0)
    FROM queues
    WHERE business_id = ? AND status = 'called' AND completed_at IS NOT NULL
`, businessId).Scan(&avgServiceTime).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate average time", "details": err.Error()})
		return
	}

	if avgServiceTime <= 0 {
		avgServiceTime = 5 // Fallback
	}

	if err := tx.Exec(`
        UPDATE queues
        SET 
            position = position - 1,
            estimated_wait_time = CAST(((position - 1) * ?) AS INTEGER)
        WHERE business_id = ?
          AND position > ?
          AND status IN ('waiting', 'called')
    `, avgServiceTime, businessId, entry.Position).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shift positions"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit failed"})
		return
	}

	// Broadcast update (if you use websockets / SSE)
	bc.Broadcaster.Subscribe(uint(businessId))

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Called Next successfully",
	})

}

//CompleteItem

//LeaveQueue

//GetQueue
