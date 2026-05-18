package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/service"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(ds *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: ds}
}

// @Summary Register Device
// @Description Register a device for push notifications
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.RegisterDeviceRequest true "Device info"
// @Success 201 {object} model.Device
// @Router /devices [post]
func (h *DeviceHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Obtener userID del contexto (del auth middleware)
	userID := r.Context().Value("user_id")
	var userIDPtr *string
	if userID != nil {
		id := userID.(string)
		userIDPtr = &id
	}

	device, err := h.deviceService.RegisterDevice(userIDPtr, req.Token, req.Platform)
	if err != nil {
		http.Error(w, "Error al registrar device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

// @Summary Get User Devices
// @Description Get all devices for current user
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []model.Device
// @Router /devices [get]
func (h *DeviceHandler) GetUserDevices(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	devices, err := h.deviceService.GetUserDevices(userID.(string))
	if err != nil {
		http.Error(w, "Error al obtener devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// @Summary Delete Device
// @Description Delete a device
// @Tags Devices
// @Security BearerAuth
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]string
// @Router /devices/{id} [delete]
func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID requerido", http.StatusBadRequest)
		return
	}

	if err := h.deviceService.DeleteDevice(id); err != nil {
		http.Error(w, "Error al eliminar device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Device eliminado"})
}