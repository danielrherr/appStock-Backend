package model

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleUsuario UserRole = "usuario"
)

type Usuario struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Nombre       *string   `json:"nombre"`
	Rol          UserRole  `json:"rol"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Nombre   *string `json:"nombre"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *Usuario `json:"user"`
	Token string   `json:"token"`
}