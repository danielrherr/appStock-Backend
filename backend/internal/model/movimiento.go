package model

import (
	"time"
)

type TipoMovimiento string

const (
	Entrada TipoMovimiento = "entrada"
	Salida  TipoMovimiento = "salida"
)

type Movimiento struct {
	ID             string         `json:"id"`
	ProductoID     string         `json:"producto_id"`
	ProductoNombre string         `json:"producto_nombre,omitempty"`
	Tipo           TipoMovimiento `json:"tipo"`
	Cantidad       int            `json:"cantidad"`
	Motivo         *string        `json:"motivo"`
	UsuarioID      *string        `json:"usuario_id"`
	Fecha          time.Time      `json:"fecha"`
}

type CreateMovimientoRequest struct {
	ProductoID string         `json:"producto_id" validate:"required"`
	Tipo      TipoMovimiento `json:"tipo" validate:"required,oneof=entrada salida"`
	Cantidad  int            `json:"cantidad" validate:"required,min=1"`
	Motivo    *string        `json:"motivo"`
}

type MovimientoFilter struct {
	ProductoID  string
	Tipo        string
	FechaDesde  string
	FechaHasta  string
	Page        int
	Limit       int
}