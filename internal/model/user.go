package model

import (
	"time"

	"github.com/google/uuid"
)

// user roles
type UserRole string

const (
	Admin     UserRole = "admin"
	Moderator UserRole = "moderator"
	AppUser   UserRole = "user"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,gte=5"`
	Password string `json:"password" binding:"required,gte=8"`
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8"`
}

type LoginUserRes struct {
	User        User   `json:"user"`
	AccessToken string `json:"accessToken"`
}
