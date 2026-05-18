package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stockapp/backend/internal/repository"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

// @Summary Get Stock by Category
// @Description Get stock grouped by category
// @Tags Reportes
// @Produce json
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /reportes/stock-categoria [get]
func (h *ReportHandler) GetStockPorCategoria(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	reporte, err := repository.GetStockPorCategoria(userID.(string))
	if err != nil {
		http.Error(w, "Error al generar reporte", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, reporte)
}

// @Summary Get Movimientos by Date Range
// @Description Get movements in a date range
// @Tags Reportes
// @Produce json
// @Security BearerAuth
// @Param desde query string false "Start date (YYYY-MM-DD)"
// @Param hasta query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} interface{}
// @Router /reportes/movimientos-fecha [get]
func (h *ReportHandler) GetMovimientosPorFecha(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	desde := r.URL.Query().Get("desde")
	hasta := r.URL.Query().Get("hasta")

	var desdeTime, hastaTime time.Time
	var err error

	if desde != "" {
		desdeTime, err = time.Parse("2006-01-02", desde)
		if err != nil {
			http.Error(w, "Formato de fecha inválido (desde)", http.StatusBadRequest)
			return
		}
	} else {
		desdeTime = time.Now().AddDate(0, -1, 0) // Default: último mes
	}

	if hasta != "" {
		hastaTime, err = time.Parse("2006-01-02", hasta)
		if err != nil {
			http.Error(w, "Formato de fecha inválido (hasta)", http.StatusBadRequest)
			return
		}
	} else {
		hastaTime = time.Now()
	}

	reporte, err := repository.GetMovimientosPorFecha(userID.(string), desdeTime, hastaTime)
	if err != nil {
		http.Error(w, "Error al generar reporte", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, reporte)
}

// @Summary Get Advanced Dashboard
// @Description Get advanced dashboard with detailed stats
// @Tags Reportes
// @Produce json
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /reportes/dashboard-avanzado [get]
func (h *ReportHandler) GetDashboardAvanzado(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	dashboard, err := repository.GetDashboardAvanzado(userID.(string))
	if err != nil {
		http.Error(w, "Error al generar dashboard", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, dashboard)
}

// Router returns the router for reports
func (h *ReportHandler) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/stock-categoria", h.GetStockPorCategoria)
		r.Get("/movimientos-fecha", h.GetMovimientosPorFecha)
		r.Get("/dashboard-avanzado", h.GetDashboardAvanzado)
	}
}