package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/service"
)

type CategoriaHandler struct {
	categoriaService *service.CategoriaService
}

func NewCategoriaHandler(categoriaService *service.CategoriaService) *CategoriaHandler {
	return &CategoriaHandler{categoriaService: categoriaService}
}

// @Summary Create Categoria
// @Description Create a new category
// @Tags Categorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateCategoriaRequest true "Categoria"
// @Success 201 {object} model.Categoria
// @Failure 400 {string} string
// @Router /categorias [post]
func (h *CategoriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	categoria, err := h.categoriaService.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoria)
}

// @Summary Get All Categorias
// @Description Get all categories
// @Tags Categorias
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /categorias [get]
func (h *CategoriaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categorias, err := h.categoriaService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": categorias})
}

// @Summary Get Categoria by ID
// @Description Get a category by ID
// @Tags Categorias
// @Produce json
// @Security BearerAuth
// @Param id path string true "Categoria ID"
// @Success 200 {object} model.Categoria
// @Router /categorias/{id} [get]
func (h *CategoriaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	categoria, err := h.categoriaService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

// @Summary Update Categoria
// @Description Update a category
// @Tags Categorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Categoria ID"
// @Param request body model.UpdateCategoriaRequest true "Categoria"
// @Success 200 {object} model.Categoria
// @Router /categorias/{id} [put]
func (h *CategoriaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateCategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	categoria, err := h.categoriaService.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

// @Summary Delete Categoria
// @Description Delete a category
// @Tags Categorias
// @Security BearerAuth
// @Param id path string true "Categoria ID"
// @Success 204 {string} string
// @Router /categorias/{id} [delete]
func (h *CategoriaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.categoriaService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}