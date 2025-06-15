package models

import (
	"time"
)

// Game represents a game event in the system
type Game struct {
	ID                  string    `json:"id,omitempty"`
	EventName           string    `json:"event_name" binding:"required"`
	StartTime           time.Time `json:"start_time" binding:"required"`
	EndTime             time.Time `json:"end_time" binding:"required"`
	Location            string    `json:"location" binding:"required"`
	CostPerPerson       float64   `json:"cost_per_person" binding:"required"`
	PlayerRequirement   int       `json:"player_requirement" binding:"required"`
	CurrentParticipants int       `json:"current_participants"`
	CreatorID           string    `json:"creator_id" binding:"required"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}

// GameCreationRequest represents the request to create a new game
type GameCreationRequest struct {
	EventName     string  `json:"event_name" binding:"required"`
	StartTime     string  `json:"start_time" binding:"required"` // Format: YYYY-MM-DDThh:mm:ss
	EndTime       string  `json:"end_time" binding:"required"`   // Format: YYYY-MM-DDThh:mm:ss
	Location      string  `json:"location" binding:"required"`
	CostPerPerson float64 `json:"cost_per_person" binding:"required"`
	PlayerRequirement   int     `json:"player_requirement" binding:"required"`
	// CreatorID comes from the JWT token
}

// GameResponse represents the response after game creation or retrieval
type GameResponse struct {
	ID                  string    `json:"id"`
	EventName           string    `json:"event_name"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	CostPerPerson       float64   `json:"cost_per_person"`
	PlayerRequirement   int       `json:"player_requirement"`
	CurrentParticipants int       `json:"current_participants"`
	CreatorID           string    `json:"creator_id"`
	CreatedAt           time.Time `json:"created_at"`
}

// GameListResponse represents a list of games
type GameListResponse struct {
	Games []GameResponse `json:"games"`
}

// JoinGameRequest represents a request to join a game
type JoinGameRequest struct {
	GameID string `json:"game_id" binding:"required"`
	// UserID comes from the JWT token
}
