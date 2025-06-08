package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {
	r := gin.Default()

	// Basic health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Get all users
	r.GET("/users", getUsers)

	// Get user by ID
	r.GET("/users/:id", getUserByID)

	// Create a new user
	r.POST("/users", createUser)

	// Start the server
	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, you would generate a unique ID and save to a database
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}
