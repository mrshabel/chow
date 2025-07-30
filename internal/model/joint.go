package model

import (
	"time"

	"github.com/google/uuid"
)

type Joint struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Distance    *float64  `json:"distance,omitempty"`
	Description *string   `json:"description"`
	IsApproved  bool      `json:"isApproved"`
	CreatorID   uuid.UUID `json:"creatorId"`
	PhotoURL    *string   `json:"photoUrl"`
	UpVotes     int       `json:"upvotes"`
	DownVotes   int       `json:"downvotes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateJoinReq struct {
	Name        string  `json:"email" binding:"required,ge=3"`
	Longitude   float64 `json:"longitude" binding:"required,longitude"`
	Latitude    float64 `json:"latitude" binding:"required,latitude"`
	Description *string `json:"description"`
}
