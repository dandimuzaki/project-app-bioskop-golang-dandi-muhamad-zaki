package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type AuthResponse struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	Token uuid.UUID `json:"token"`
}