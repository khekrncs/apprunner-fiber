package models

import (
	"time"
)

type User struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string            `json:"email" validate:"required,email"`
	Name     string            `json:"name" validate:"required"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type UpdateUserRequest struct {
	Email    string            `json:"email,omitempty" validate:"omitempty,email"`
	Name     string            `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
