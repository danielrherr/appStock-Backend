package model

import (
	"time"
)

type Categoria struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion *string   `json:"descripcion"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoriaRequest struct {
	Nombre      string  `json:"nombre" validate:"required"`
	Descripcion *string `json:"descripcion"`
}

type UpdateCategoriaRequest struct {
	Nombre      *string `json:"nombre"`
	Descripcion *string `json:"descripcion"`
}