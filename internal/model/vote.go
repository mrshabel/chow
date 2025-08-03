package model

import (
	"time"

	"github.com/google/uuid"
)

type VoteDirection string

const (
	UpVote   VoteDirection = "up"
	DownVote VoteDirection = "down"
)

type Vote struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"userId"`
	JointID   uuid.UUID     `json:"jointId"`
	Direction VoteDirection `json:"direction"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type JointVoter struct {
	UserID    uuid.UUID     `json:"userId"`
	Username  string        `json:"username"`
	JointID   uuid.UUID     `json:"jointId"`
	Direction VoteDirection `json:"direction"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type VoteJointReq struct {
	Direction VoteDirection `json:"direction" binding:"required,oneof=up down"`
}
