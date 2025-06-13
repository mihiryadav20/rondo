package routes

import (
	"github.com/gin-gonic/gin"
	
	"rondo/handlers"
	"rondo/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine) {
	// Auth routes
	auth := r.Group("/auth")
	{
		// Public OTP endpoints
		otp := auth.Group("/otp")
		{
			otp.POST("/request", handlers.RequestOTP)
			otp.POST("/verify", handlers.VerifyOTP)
		}
		
		// Protected routes - require JWT authentication
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User registration - requires JWT from OTP verification
			protected.POST("/register", handlers.RegisterUser)
		}
	}
	
	// User routes - protected by JWT authentication
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware()) // Apply JWT middleware to all user routes
	{
		users.GET("/:phone", handlers.GetUserProfile)
	}
}
