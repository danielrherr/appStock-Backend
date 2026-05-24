package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/repository"
)

type ProductoService struct {
	uploadDir string
}

func NewProductoService(uploadDir string) *ProductoService {
	return &ProductoService{uploadDir: uploadDir}
}

func (s *ProductoService) CreateProducto(req *model.CreateProductoRequest) (*model.Producto, error) {
	// Check if codigo exists
	exists, err := repository.CodigoExists(req.Codigo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ya existe un producto con ese código")
	}

	return repository.CreateProducto(req)
}

func (s *ProductoService) GetByID(id string) (*model.Producto, error) {
	return repository.GetProductoByID(id)
}

func (s *ProductoService) GetByCodigo(codigo string) (*model.Producto, error) {
	producto, err := repository.GetProductoByCodigo(codigo)
	if err != nil {
		return nil, errors.New("producto no encontrado")
	}
	return producto, nil
}

func (s *ProductoService) GetAll(page, limit int, search, categoriaID string, stockBajo bool) ([]model.Producto, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return repository.GetProductos(page, limit, search, categoriaID, stockBajo)
}

func (s *ProductoService) GetStockBajo() ([]model.Producto, error) {
	return repository.GetProductosStockBajo()
}

func (s *ProductoService) Update(id string, req *model.UpdateProductoRequest) (*model.Producto, error) {
	// Check if exists
	_, err := repository.GetProductoByID(id)
	if err != nil {
		return nil, errors.New("producto no encontrado")
	}

	// Check codigo unique
	if req.Codigo != nil {
		exists, err := repository.CodigoExistsForOtherProduct(*req.Codigo, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("ya existe un producto con ese código")
		}
	}

	return repository.UpdateProducto(id, req)
}

func (s *ProductoService) UploadImagen(productoID string, filename string, file io.Reader) (string, error) {
	// Check if producto exists
	producto, err := repository.GetProductoByID(productoID)
	if err != nil {
		return "", errors.New("producto no encontrado")
	}

	// Create upload dir if not exists
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}

	// Generate filename
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", producto.ID, ext)
	filePath := filepath.Join(s.uploadDir, newFilename)

	// Save file
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return "", err
	}

	// Update producto
	imagenPath := "/uploads/" + newFilename
	updated, err := repository.UpdateImagen(productoID, imagenPath)
	if err != nil {
		return "", err
	}

	if updated.Imagen != nil {
		return *updated.Imagen, nil
	}
	return "", nil
}

func (s *ProductoService) Delete(id string) error {
	// Check if exists
	_, err := repository.GetProductoByID(id)
	if err != nil {
		return errors.New("producto no encontrado")
	}

	return repository.DeleteProducto(id)
}
