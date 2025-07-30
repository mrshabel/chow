package model

import (
	"time"

	"github.com/google/uuid"
)

type ComplaintStatus string

const (
	Open     ComplaintStatus = "open"
	Resolved ComplaintStatus = "resolved"
)

type Complaint struct {
	ID        uuid.UUID       `json:"id"`
	JointID   uuid.UUID       `json:"jointId"`
	UserID    uuid.UUID       `json:"userId"`
	Reason    string          `json:"reason"`
	Status    ComplaintStatus `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type CreateComplaintReq struct {
	Reason string `json:"reason" binding:"required,gt=5"`
}
