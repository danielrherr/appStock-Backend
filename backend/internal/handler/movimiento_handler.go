package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/middleware"
	"github.com/stockapp/backend/internal/service"
)

type MovimientoHandler struct {
	movimientoService *service.MovimientoService
}

func NewMovimientoHandler(movimientoService *service.MovimientoService) *MovimientoHandler {
	return &MovimientoHandler{movimientoService: movimientoService}
}

// @Summary Create Movimiento
// @Description Create a stock movement (entrada/salida)
// @Tags Movimientos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateMovimientoRequest true "Movimiento"
// @Success 201 {object} model.Movimiento
// @Failure 400 {string} string
// @Router /movimientos [post]
func (h *MovimientoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateMovimientoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserID(r)
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	movimiento, err := h.movimientoService.Create(req.ProductoID, &req, userIDPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movimiento)
}

// @Summary Get All Movimientos
// @Description Get all stock movements with filters
// @Tags Movimientos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param producto_id query string false "Filter by product"
// @Param tipo query string false "Filter by type (entrada/salida)"
// @Param fecha_desde query string false "Filter from date"
// @Param fecha_hasta query string false "Filter to date"
// @Success 200 {object} map[string]interface{}
// @Router /movimientos [get]
func (h *MovimientoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter := model.MovimientoFilter{
		ProductoID:  r.URL.Query().Get("producto_id"),
		Tipo:       r.URL.Query().Get("tipo"),
		FechaDesde: r.URL.Query().Get("fecha_desde"),
		FechaHasta: r.URL.Query().Get("fecha_hasta"),
		Page:       getIntParam(r, "page", 1),
		Limit:      getIntParam(r, "limit", 20),
	}

	movimientos, total, err := h.movimientoService.GetAll(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": movimientos,
		"pagination": map[string]int{
			"page":  filter.Page,
			"limit": filter.Limit,
			"total": total,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// @Summary Get Movimientos by Producto
// @Description Get all movements for a specific product
// @Tags Movimientos
// @Produce json
// @Security BearerAuth
// @Param producto_id path string true "Producto ID"
// @Success 200 {object} []model.Movimiento
// @Router /movimientos/producto/{producto_id} [get]
func (h *MovimientoHandler) GetByProducto(w http.ResponseWriter, r *http.Request) {
	productoID := chi.URLParam(r, "producto_id")

	movimientos, err := h.movimientoService.GetByProducto(productoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": movimientos})
}

// @Summary Get Dashboard
// @Description Get dashboard statistics
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /dashboard [get]
func (h *MovimientoHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.movimientoService.GetDashboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}