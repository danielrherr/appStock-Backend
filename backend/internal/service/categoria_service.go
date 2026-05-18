package service

import (
	"errors"

	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/repository"
)

type CategoriaService struct{}

func NewCategoriaService() *CategoriaService {
	return &CategoriaService{}
}

func (s *CategoriaService) Create(req *model.CreateCategoriaRequest) (*model.Categoria, error) {
	// Check if nombre exists
	exists, err := repository.CategoriaNameExists(req.Nombre)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ya existe una categoría con ese nombre")
	}

	return repository.CreateCategoria(req.Nombre, req.Descripcion)
}

func (s *CategoriaService) GetByID(id string) (*model.Categoria, error) {
	return repository.GetCategoriaByID(id)
}

func (s *CategoriaService) GetAll() ([]model.Categoria, error) {
	return repository.GetAllCategorias()
}

func (s *CategoriaService) Update(id string, req *model.UpdateCategoriaRequest) (*model.Categoria, error) {
	// Check if exists
	cat, err := repository.GetCategoriaByID(id)
	if err != nil {
		return nil, errors.New("categoría no encontrada")
	}

	// Check if nombre is being changed and if it already exists
	if req.Nombre != nil && *req.Nombre != cat.Nombre {
		exists, err := repository.CategoriaNameExists(*req.Nombre)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("ya existe una categoría con ese nombre")
		}
	}

	nombre := req.Nombre
	if nombre == nil {
		nombre = &cat.Nombre
	}

	descripcion := req.Descripcion
	if descripcion == nil {
		descripcion = cat.Descripcion
	}

	return repository.UpdateCategoria(id, *nombre, descripcion)
}

func (s *CategoriaService) Delete(id string) error {
	// Check if exists
	_, err := repository.GetCategoriaByID(id)
	if err != nil {
		return errors.New("categoría no encontrada")
	}

	return repository.DeleteCategoria(id)
}