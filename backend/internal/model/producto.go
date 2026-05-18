package model

import (
	"time"
)

type Producto struct {
	ID            string    `json:"id"`
	Codigo        string    `json:"codigo"`
	CodigoBarras  *string   `json:"codigo_barras"`
	Nombre        string    `json:"nombre"`
	Descripcion   *string   `json:"descripcion"`
	CategoriaID   *string   `json:"categoria_id"`
	CategoriaNombre string  `json:"categoria_nombre,omitempty"`
	Precio        float64   `json:"precio"`
	StockActual   int       `json:"stock_actual"`
	StockMinimo   int       `json:"stock_minimo"`
	Imagen        *string   `json:"imagen"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateProductoRequest struct {
	Codigo        string  `json:"codigo" validate:"required"`
	CodigoBarras  *string `json:"codigo_barras"`
	Nombre        string  `json:"nombre" validate:"required"`
	Descripcion   *string `json:"descripcion"`
	CategoriaID   *string `json:"categoria_id"`
	Precio        float64 `json:"precio" validate:"required"`
	StockActual   int     `json:"stock_actual"`
	StockMinimo   int     `json:"stock_minimo"`
}

type UpdateProductoRequest struct {
	Codigo        *string `json:"codigo"`
	CodigoBarras  *string `json:"codigo_barras"`
	Nombre        *string `json:"nombre"`
	Descripcion   *string `json:"descripcion"`
	CategoriaID   *string `json:"categoria_id"`
	Precio        *float64 `json:"precio"`
	StockMinimo   *int    `json:"stock_minimo"`
}