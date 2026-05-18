package repository

import (
	"fmt"
	"time"

	"github.com/stockapp/backend/internal/model"
)

type StockPorCategoria struct {
	CategoriaID   string  `json:"categoria_id"`
	CategoriaName string  `json:"categoria_nombre"`
	TotalProductos int     `json:"total_productos"`
	TotalStock    int     `json:"total_stock"`
	TotalValor    float64 `json:"total_valor"`
}

func GetStockPorCategoria(userID string) ([]StockPorCategoria, error) {
	rows, err := DB.Query(`
		SELECT 
			c.id as categoria_id,
			c.nombre as categoria_nombre,
			COUNT(p.id) as total_productos,
			COALESCE(SUM(p.stock), 0) as total_stock,
			COALESCE(SUM(p.stock * p.precio), 0) as total_valor
		FROM categorias c
		LEFT JOIN productos p ON p.categoria_id = c.id AND p.activo = 1
		WHERE c.usuario_id = ?
		GROUP BY c.id, c.nombre
		ORDER BY total_stock DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []StockPorCategoria
	for rows.Next() {
		var r StockPorCategoria
		if err := rows.Scan(&r.CategoriaID, &r.CategoriaName, &r.TotalProductos, &r.TotalStock, &r.TotalValor); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

type MovimientoPorFecha struct {
	Fecha            string  `json:"fecha"`
	TotalEntradas    int     `json:"total_entradas"`
	TotalSalidas     int     `json:"total_salidas"`
	EntradasMonto    float64 `json:"entradas_monto"`
	SalidasMonto     float64 `json:"salidas_monto"`
	MovimientosCount int     `json:"movimientos_count"`
}

func GetMovimientosPorFecha(userID string, desde, hasta time.Time) ([]MovimientoPorFecha, error) {
	rows, err := DB.Query(`
		SELECT 
			DATE(m.created_at) as fecha,
			SUM(CASE WHEN m.tipo = 'entrada' THEN m.cantidad ELSE 0 END) as total_entradas,
			SUM(CASE WHEN m.tipo = 'salida' THEN m.cantidad ELSE 0 END) as total_salidas,
			SUM(CASE WHEN m.tipo = 'entrada' THEN m.cantidad * p.precio ELSE 0 END) as entradas_monto,
			SUM(CASE WHEN m.tipo = 'salida' THEN m.cantidad * p.precio ELSE 0 END) as salidas_monto,
			COUNT(m.id) as movimientos_count
		FROM movimientos m
		JOIN productos p ON p.id = m.producto_id
		WHERE p.usuario_id = ? AND DATE(m.created_at) BETWEEN ? AND ?
		GROUP BY DATE(m.created_at)
		ORDER BY fecha DESC
	`, userID, desde.Format("2006-01-02"), hasta.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MovimientoPorFecha
	for rows.Next() {
		var r MovimientoPorFecha
		if err := rows.Scan(&r.Fecha, &r.TotalEntradas, &r.TotalSalidas, &r.EntradasMonto, &r.SalidasMonto, &r.MovimientosCount); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

type DashboardAvanzado struct {
	TotalProductos     int     `json:"total_productos"`
	TotalCategorias    int     `json:"total_categorias"`
	TotalStock         int     `json:"total_stock"`
	TotalValor         float64 `json:"total_valor"`
	ProductosStockBajo int     `json:"productos_stock_bajo"`
	MovimientosHoy     int     `json:"movimientos_hoy"`
	EntradasMes        int     `json:"entradas_mes"`
	SalidasMes         int     `json:"salidas_mes"`
	TopProductos       []model.Producto `json:"top_productos"`
}

func GetDashboardAvanzado(userID string) (*DashboardAvanzado, error) {
	dash := &DashboardAvanzado{}

	// Totales
	err := DB.QueryRow(`
		SELECT 
			COUNT(*) as total_productos,
			COALESCE(SUM(stock), 0) as total_stock,
			COALESCE(SUM(stock * precio), 0) as total_valor
		FROM productos 
		WHERE usuario_id = ? AND activo = 1
	`, userID).Scan(&dash.TotalProductos, &dash.TotalStock, &dash.TotalValor)
	if err != nil {
		return nil, err
	}

	// Total categorías
	err = DB.QueryRow(`SELECT COUNT(*) FROM categorias WHERE usuario_id = ?`, userID).Scan(&dash.TotalCategorias)
	if err != nil {
		return nil, err
	}

	// Stock bajo
	err = DB.QueryRow(`
		SELECT COUNT(*) FROM productos 
		WHERE usuario_id = ? AND stock <= stock_minimo AND activo = 1
	`, userID).Scan(&dash.ProductosStockBajo)
	if err != nil {
		return nil, err
	}

	// Movimientos hoy
	fechaHoy := time.Now().Format("2006-01-02")
	err = DB.QueryRow(`
		SELECT COUNT(*) FROM movimientos m
		JOIN productos p ON p.id = m.producto_id
		WHERE p.usuario_id = ? AND DATE(m.created_at) = ?
	`, userID, fechaHoy).Scan(&dash.MovimientosHoy)
	if err != nil {
		return nil, err
	}

	// Entradas/salidas del mes actual
	mesActual := time.Now().Format("2006-01")
	err = DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN m.tipo = 'entrada' THEN m.cantidad ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN m.tipo = 'salida' THEN m.cantidad ELSE 0 END), 0)
		FROM movimientos m
		JOIN productos p ON p.id = m.producto_id
		WHERE p.usuario_id = ? AND strftime('%Y-%m', m.created_at) = ?
	`, userID, mesActual).Scan(&dash.EntradasMes, &dash.SalidasMes)
	if err != nil {
		return nil, err
	}

	// Top productos (por movimiento)
	rows, err := DB.Query(`
		SELECT p.id, p.nombre, p.stock, p.precio
		FROM productos p
		LEFT JOIN movimientos m ON m.producto_id = p.id
		WHERE p.usuario_id = ? AND p.activo = 1
		GROUP BY p.id
		ORDER BY COUNT(m.id) DESC
		LIMIT 5
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prod model.Producto
		if err := rows.Scan(&prod.ID, &prod.Nombre, &prod.StockActual, &prod.Precio); err != nil {
			return nil, err
		}
		dash.TopProductos = append(dash.TopProductos, prod)
	}

	return dash, nil
}

// Helper to avoid import issues
func init() {
	_ = fmt.Sprintf("%s", time.Now())
}