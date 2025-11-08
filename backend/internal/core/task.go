package core

import (
	"errors"
	"strings"
	"time"
)

// Task is the domain entity used across the backend.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"` // todo | doing | done
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewTask is the input DTO for creating tasks.
type NewTask struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

// UpdateTask is the input DTO for updates (all fields optional).
type UpdateTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

var (
	ErrInvalidTitle  = errors.New("invalid_title")
	ErrInvalidDesc   = errors.New("invalid_description")
	ErrInvalidStatus = errors.New("invalid_status")
)

func ValidateNew(in NewTask) error {
	title := strings.TrimSpace(in.Title)
	if len(title) == 0 || len(title) > 140 {
		return ErrInvalidTitle
	}
	if in.Description != nil && len(*in.Description) > 1000 {
		return ErrInvalidDesc
	}
	return nil
}

func ValidateUpdate(in UpdateTask) error {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if len(title) == 0 || len(title) > 140 {
			return ErrInvalidTitle
		}
	}
	if in.Description != nil && len(*in.Description) > 1000 {
		return ErrInvalidDesc
	}
	if in.Status != nil {
		switch *in.Status {
		case "todo", "doing", "done":
		// ok
		default:
			return ErrInvalidStatus
		}
	}
	return nil
}
