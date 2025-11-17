package routes

import (
	"smart-queue/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, bc *controller.BusinessController) {

	api := r.Group("/api")

	{
		api.GET("/businesses", bc.ListBusiness)
		api.GET("/businesses/:id", bc.FetchBusinessDetailById)
		// api.GET("/businesses/events/:businessID", bc.StreamSSE)
		// api.POST("/businesses/broadcast/:businessID", bc.BroadcastMessage)
		api.POST("/businesses/join_queue/:business_id", bc.JoinQueue)
		api.POST("/businesses/leave_queue/:business_id", bc.LeaveQueue)
		api.POST("/businesses/call_next/:business_id", bc.CallNext)
	}

}
