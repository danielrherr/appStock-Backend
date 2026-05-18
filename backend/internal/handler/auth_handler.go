package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Register
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Register request"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {string} string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// @Summary Login
// @Description Authenticate user and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} model.AuthResponse
// @Failure 401 {string} string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}