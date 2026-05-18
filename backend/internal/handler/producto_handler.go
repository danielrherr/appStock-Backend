package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/service"
)

type ProductoHandler struct {
	productoService *service.ProductoService
}

func NewProductoHandler(productoService *service.ProductoService) *ProductoHandler {
	return &ProductoHandler{productoService: productoService}
}

// @Summary Create Producto
// @Description Create a new product
// @Tags Productos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateProductoRequest true "Producto"
// @Success 201 {object} model.Producto
// @Failure 400 {string} string
// @Router /productos [post]
func (h *ProductoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProductoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	producto, err := h.productoService.CreateProducto(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(producto)
}

// @Summary Get All Productos
// @Description Get all products with pagination and filters
// @Tags Productos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or code"
// @Param categoria_id query string false "Filter by category"
// @Param stock_bajo query bool false "Filter by low stock"
// @Success 200 {object} map[string]interface{}
// @Router /productos [get]
func (h *ProductoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page := getIntParam(r, "page", 1)
	limit := getIntParam(r, "limit", 20)
	search := r.URL.Query().Get("search")
	categoriaID := r.URL.Query().Get("categoria_id")
	stockBajo := r.URL.Query().Get("stock_bajo") == "true"

	productos, total, err := h.productoService.GetAll(page, limit, search, categoriaID, stockBajo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": productos,
		"pagination": map[string]int{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// @Summary Get Producto by ID
// @Description Get a product by ID
// @Tags Productos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Producto ID"
// @Success 200 {object} model.Producto
// @Router /productos/{id} [get]
func (h *ProductoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	producto, err := h.productoService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// @Summary Get Producto by Barcode
// @Description Get a product by barcode code
// @Tags Productos
// @Produce json
// @Security BearerAuth
// @Param codigo path string true "Barcode code"
// @Success 200 {object} model.Producto
// @Router /productos/barcode/{codigo} [get]
func (h *ProductoHandler) GetByBarcode(w http.ResponseWriter, r *http.Request) {
	codigo := chi.URLParam(r, "codigo")

	producto, err := h.productoService.GetByCodigo(codigo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// @Summary Get Low Stock Products
// @Description Get all products with stock below minimum
// @Tags Productos
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /productos/stock-bajo [get]
func (h *ProductoHandler) GetStockBajo(w http.ResponseWriter, r *http.Request) {
	productos, err := h.productoService.GetStockBajo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": productos})
}

// @Summary Update Producto
// @Description Update a product
// @Tags Productos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Producto ID"
// @Param request body model.UpdateProductoRequest true "Producto"
// @Success 200 {object} model.Producto
// @Router /productos/{id} [put]
func (h *ProductoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateProductoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	producto, err := h.productoService.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// @Summary Upload Product Image
// @Description Upload an image for a product
// @Tags Productos
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Producto ID"
// @Param imagen formData file true "Product image"
// @Success 200 {object} map[string]string
// @Router /productos/{id}/imagen [post]
func (h *ProductoHandler) UploadImagen(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	file, header, err := r.FormFile("imagen")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imagen, err := h.productoService.UploadImagen(id, header.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"imagen": imagen})
}

// @Summary Delete Producto
// @Description Delete a product
// @Tags Productos
// @Security BearerAuth
// @Param id path string true "Producto ID"
// @Success 204 {string} string
// @Router /productos/{id} [delete]
func (h *ProductoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.productoService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getIntParam(r *http.Request, name string, defaultValue int) int {
	var value int
	fmt.Sscanf(r.URL.Query().Get(name), "%d", &value)
	if value == 0 {
		return defaultValue
	}
	return value
}