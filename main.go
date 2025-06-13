package main

import (
	"github.com/gin-gonic/gin"
	
	"rondo/config"
	"rondo/handlers"
	"rondo/routes"
	"rondo/utils"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Initialize Twilio client
	twilioClient := utils.InitTwilio()
	
	// Initialize handlers
	handlers.InitHandlers(twilioClient)

	// Setup router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	r.Run(":8080")
}
