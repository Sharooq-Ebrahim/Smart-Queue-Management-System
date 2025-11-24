package routes

import (
	"smart-queue/controller"
	"smart-queue/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, bc *controller.BusinessController, auth *controller.AuthController) {

	api := r.Group("/api")

	{

		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)
		authorized := r.Group("/api")

		authorized.Use(middleware.AuthMiddlware())

		authorized.GET("/businesses", bc.ListBusiness)
		authorized.GET("/businesses/:id", bc.FetchBusinessDetailById)
		authorized.GET("/businesses/events/:businessID", bc.StreamSSE)
		authorized.POST("/businesses/join_queue/:business_id", bc.JoinQueue)
		authorized.POST("/businesses/leave_queue/:business_id", bc.LeaveQueue)
		authorized.POST("/businesses/call_next/:business_id", bc.CallNext)
	}

}
