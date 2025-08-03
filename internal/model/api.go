package model

import "github.com/google/uuid"

const (
	DefaultPageSize = 10
	DefaultPage     = 1
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"pageSize" binding:"omitempty,gte=1,lte=1000"`
}

func (p *PaginationQuery) GetOffsetAndLimit() (int, int) {
	if p.Page == 0 {
		p.Page = DefaultPage
	}
	if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}

	offset := (p.Page - 1) * p.PageSize
	// offset and limit
	return offset, p.PageSize
}

// params with only ID
type IDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetID returns a uuid representation of the ID param string
func (p *IDParam) GetID() uuid.UUID {
	id, _ := uuid.Parse(p.ID)
	return id
}
