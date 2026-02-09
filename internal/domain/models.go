package domain

import (
	"time"

	"github.com/google/uuid"
)

// @Description FAQ is an internal model used by service and repository.
type FAQ struct {
	ID        uuid.UUID
	Title     string
	Content   string
	Position  int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// @Description CreateFAQRequest describes request body for creating a FAQ.
type CreateFAQRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position int    `json:"position"`
	IsActive *bool  `json:"is_active"`
}

// @Description UpdateFAQRequest describes request body for updating a FAQ.
type UpdateFAQRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position int    `json:"position"`
	IsActive *bool  `json:"is_active"`
}

type CreateFAQInput struct {
	Title    string
	Content  string
	Position int
	IsActive bool
}

type UpdateFAQInput struct {
	Title    string
	Content  string
	Position int
	IsActive bool
}

// @Description FAQListItemResponse is a short FAQ representation used in lists.
type FAQListItemResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Position int       `json:"position"`
}

// @Description FAQFullResponse is a full FAQ representation.
type FAQFullResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Position int       `json:"position"`
	IsActive bool      `json:"is_active"`
}

// @Description DataResponse wraps API response payloads.
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// @Description FAQListResponse wraps a list response.
type FAQListResponse struct {
	Data []FAQListItemResponse `json:"data"`
}

// @Description FAQResponse wraps a single FAQ response.
type FAQResponse struct {
	Data FAQFullResponse `json:"data"`
}

// @Description MessageResponse is a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// @Description ErrorResponse describes an error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}
