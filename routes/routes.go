package routes

import (
	"github.com/gin-gonic/gin"
	
	"rondo/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine) {
	// Auth routes
	auth := r.Group("/auth")
	{
		otp := auth.Group("/otp")
		{
			otp.POST("/request", handlers.RequestOTP)
			otp.POST("/verify", handlers.VerifyOTP)
		}
	}
}
