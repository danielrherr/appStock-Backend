package service

import (
	"errors"

	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/repository"
)

type MovimientoService struct{}

func NewMovimientoService() *MovimientoService {
	return &MovimientoService{}
}

func (s *MovimientoService) Create(productoID string, req *model.CreateMovimientoRequest, usuarioID *string) (*model.Movimiento, error) {
	// Check if producto exists
	producto, err := repository.GetProductoByID(productoID)
	if err != nil {
		return nil, errors.New("producto no encontrado")
	}

	// Validate stock for salida
	if req.Tipo == model.Salida {
		if producto.StockActual < req.Cantidad {
			return nil, errors.New("stock insuficiente")
		}
	}

	return repository.CreateMovimiento(productoID, req, usuarioID)
}

func (s *MovimientoService) GetByID(id string) (*model.Movimiento, error) {
	return repository.GetMovimientoByID(id)
}

func (s *MovimientoService) GetAll(filter model.MovimientoFilter) ([]model.Movimiento, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	return repository.GetMovimientos(filter)
}

func (s *MovimientoService) GetByProducto(productoID string) ([]model.Movimiento, error) {
	// Check if producto exists
	_, err := repository.GetProductoByID(productoID)
	if err != nil {
		return nil, errors.New("producto no encontrado")
	}

	return repository.GetMovimientosByProducto(productoID)
}

type DashboardStats struct {
	TotalProductos      int     `json:"total_productos"`
	ProductosStockBajo  int     `json:"productos_stock_bajo"`
	MovimientosHoy      int     `json:"movimientos_hoy"`
	ValorTotalStock     float64 `json:"valor_total_stock"`
}

func (s *MovimientoService) GetDashboard() (*DashboardStats, error) {
	// Total productos
	productos, _, err := repository.GetProductos(1, 1, "", "", false)
	if err != nil {
		return nil, err
	}
	totalProductos := len(productos)
	if productos, _, err = repository.GetProductos(1, 100000, "", "", false); err == nil {
		totalProductos = len(productos)
	}

	// Stock bajo
	stockBajo, err := repository.GetProductosStockBajo()
	if err != nil {
		return nil, err
	}

	// Movimientos hoy
	movimientosHoy, err := repository.GetMovimientosHoy()
	if err != nil {
		return nil, err
	}

	// Valor total stock
	var valorTotal float64
	for _, p := range productos {
		valorTotal += p.Precio * float64(p.StockActual)
	}

	return &DashboardStats{
		TotalProductos:     totalProductos,
		ProductosStockBajo: len(stockBajo),
		MovimientosHoy:     movimientosHoy,
		ValorTotalStock:    valorTotal,
	}, nil
}