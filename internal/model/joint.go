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

type CreateJointReq struct {
	Name        string  `json:"name" binding:"required,gte=3"`
	Longitude   float64 `json:"longitude" binding:"required,longitude"`
	Latitude    float64 `json:"latitude" binding:"required,latitude"`
	Description *string `json:"description"`
}

type NearbyJointsQuery struct {
	Radius    float64 `form:"radius" binding:"required,gt=0,lte=5000"`
	Longitude float64 `form:"longitude" binding:"required,longitude"`
	Latitude  float64 `form:"latitude" binding:"required,latitude"`
}
