package handlers

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"rondo/models"
)

// In-memory storage for games (would be replaced with a database in production)
var games = make(map[string]models.Game)

// CreateGame handles the creation of a new game
func CreateGame(c *gin.Context) {
	// Get user ID from JWT claims
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse request body
	var req models.GameCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format. Use YYYY-MM-DDThh:mm:ssZ"})
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format. Use YYYY-MM-DDThh:mm:ssZ"})
		return
	}
	
	// Validate times
	if startTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time cannot be in the past"})
		return
	}
	
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}
	
	// Create new game
	gameID := uuid.New().String()
	now := time.Now()
	
	game := models.Game{
		ID:                  gameID,
		EventName:           req.EventName,
		StartTime:           startTime,
		EndTime:             endTime,
		Location:            req.Location,
		CostPerPerson:       req.CostPerPerson,
		PlayerRequirement:   req.PlayerRequirement,
		CurrentParticipants: 0, // Initially no participants
		CreatorID:           userID.(string),
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	
	// Save game (in a real app, this would be in a database)
	games[gameID] = game
	
	// Return response
	c.JSON(http.StatusCreated, models.GameResponse{
		ID:                  game.ID,
		EventName:           game.EventName,
		StartTime:           game.StartTime,
		EndTime:             game.EndTime,
		Location:            game.Location,
		CostPerPerson:       game.CostPerPerson,
		PlayerRequirement:   game.PlayerRequirement,
		CurrentParticipants: game.CurrentParticipants,
		CreatorID:           game.CreatorID,
		CreatedAt:           game.CreatedAt,
	})
}

// GetGame retrieves a specific game by ID
func GetGame(c *gin.Context) {
	gameID := c.Param("id")
	
	game, exists := games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	
	c.JSON(http.StatusOK, models.GameResponse{
		ID:                  game.ID,
		EventName:           game.EventName,
		StartTime:           game.StartTime,
		EndTime:             game.EndTime,
		Location:            game.Location,
		CostPerPerson:       game.CostPerPerson,
		PlayerRequirement:   game.PlayerRequirement,
		CurrentParticipants: game.CurrentParticipants,
		CreatorID:           game.CreatorID,
		CreatedAt:           game.CreatedAt,
	})
}

// ListGames returns all available games
func ListGames(c *gin.Context) {
	var gameList []models.GameResponse
	
	for _, game := range games {
		gameList = append(gameList, models.GameResponse{
			ID:                  game.ID,
			EventName:           game.EventName,
			StartTime:           game.StartTime,
			EndTime:             game.EndTime,
			Location:            game.Location,
			CostPerPerson:       game.CostPerPerson,
			PlayerRequirement:   game.PlayerRequirement,
			CurrentParticipants: game.CurrentParticipants,
			CreatorID:           game.CreatorID,
			CreatedAt:           game.CreatedAt,
		})
	}
	
	c.JSON(http.StatusOK, models.GameListResponse{
		Games: gameList,
	})
}

// JoinGame allows a user to join a game
func JoinGame(c *gin.Context) {
	// Get user ID from JWT claims
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse request body
	var req models.JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get game
	game, exists := games[req.GameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	
	// Check if game is full
	if game.CurrentParticipants >= game.PlayerRequirement {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game is already full"})
		return
	}
	
	// Check if game has already started
	if game.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game has already started"})
		return
	}
	
	// Update participant count (in a real app, we would add the user to a participants list)
	game.CurrentParticipants++
	game.UpdatedAt = time.Now()
	games[req.GameID] = game
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined the game",
		"game": models.GameResponse{
			ID:                  game.ID,
			EventName:           game.EventName,
			StartTime:           game.StartTime,
			EndTime:             game.EndTime,
			Location:            game.Location,
			CostPerPerson:       game.CostPerPerson,
			PlayerRequirement:   game.PlayerRequirement,
			CurrentParticipants: game.CurrentParticipants,
			CreatorID:           game.CreatorID,
			CreatedAt:           game.CreatedAt,
		},
	})
}
