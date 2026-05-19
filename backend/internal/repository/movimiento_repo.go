package repository

import (
	"fmt"
	"time"

	"github.com/stockapp/backend/internal/model"
)

func CreateMovimiento(productoID string, req *model.CreateMovimientoRequest, usuarioID *string) (*model.Movimiento, error) {
	// PostgreSQL generates UUID automatically via gen_random_uuid()
	var id string
	err := DB.QueryRow(
		`INSERT INTO movimientos (producto_id, tipo, cantidad, motivo, usuario_id) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		productoID, req.Tipo, req.Cantidad, req.Motivo, usuarioID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Update stock
	err = UpdateStock(productoID, req.Cantidad, req.Tipo == model.Entrada)
	if err != nil {
		return nil, err
	}

	return GetMovimientoByID(id)
}

func GetMovimientoByID(id string) (*model.Movimiento, error) {
	var m model.Movimiento
	err := DB.QueryRow(
		`SELECT m.id, m.producto_id, p.nombre, m.tipo, m.cantidad, m.motivo, m.usuario_id, m.fecha 
		 FROM movimientos m 
		 JOIN productos p ON m.producto_id = p.id 
		 WHERE m.id = $1`, id,
	).Scan(&m.ID, &m.ProductoID, &m.ProductoNombre, &m.Tipo, &m.Cantidad, &m.Motivo, &m.UsuarioID, &m.Fecha)
	
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func GetMovimientos(filter model.MovimientoFilter) ([]model.Movimiento, int, error) {
	offset := (filter.Page - 1) * filter.Limit
	
	baseQuery := `FROM movimientos m JOIN productos p ON m.producto_id = p.id WHERE 1=1`
	var args []interface{}
	argNum := 1
	
	if filter.ProductoID != "" {
		baseQuery += fmt.Sprintf(` AND m.producto_id = $%d`, argNum)
		args = append(args, filter.ProductoID)
		argNum++
	}
	
	if filter.Tipo != "" {
		baseQuery += fmt.Sprintf(` AND m.tipo = $%d`, argNum)
		args = append(args, filter.Tipo)
		argNum++
	}
	
	if filter.FechaDesde != "" {
		baseQuery += fmt.Sprintf(` AND m.fecha >= $%d`, argNum)
		if t, err := time.Parse("2006-01-02", filter.FechaDesde); err == nil {
			args = append(args, t)
			argNum++
		}
	}
	
	if filter.FechaHasta != "" {
		baseQuery += fmt.Sprintf(` AND m.fecha <= $%d`, argNum)
		if t, err := time.Parse("2006-01-02", filter.FechaHasta); err == nil {
			args = append(args, t.Add(24*time.Hour))
			argNum++
		}
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := fmt.Sprintf(`SELECT m.id, m.producto_id, p.nombre, m.tipo, m.cantidad, m.motivo, m.usuario_id, m.fecha %s ORDER BY m.fecha DESC LIMIT $%d OFFSET $%d`, baseQuery, argNum, argNum+1)
	args = append(args, filter.Limit, offset)
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movimientos []model.Movimiento
	for rows.Next() {
		var m model.Movimiento
		if err := rows.Scan(&m.ID, &m.ProductoID, &m.ProductoNombre, &m.Tipo, &m.Cantidad, &m.Motivo, &m.UsuarioID, &m.Fecha); err != nil {
			return nil, 0, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, total, nil
}

func GetMovimientosByProducto(productoID string) ([]model.Movimiento, error) {
	rows, err := DB.Query(
		`SELECT m.id, m.producto_id, p.nombre, m.tipo, m.cantidad, m.motivo, m.usuario_id, m.fecha 
		 FROM movimientos m 
		 JOIN productos p ON m.producto_id = p.id 
		 WHERE m.producto_id = $1 
		 ORDER BY m.fecha DESC`,
		productoID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimientos []model.Movimiento
	for rows.Next() {
		var m model.Movimiento
		if err := rows.Scan(&m.ID, &m.ProductoID, &m.ProductoNombre, &m.Tipo, &m.Cantidad, &m.Motivo, &m.UsuarioID, &m.Fecha); err != nil {
			return nil, err
		}
		movimientos = append(movimientos, m)
	}
	return movimientos, nil
}

func GetMovimientosHoy() (int, error) {
	var count int
	err := DB.QueryRow(
		`SELECT COUNT(*) FROM movimientos WHERE date(fecha) = CURRENT_DATE`,
	).Scan(&count)
	return count, err
}