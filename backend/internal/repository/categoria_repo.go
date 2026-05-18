package repository

import (
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/utils"
)

func CreateCategoria(nombre string, descripcion *string) (*model.Categoria, error) {
	id := utils.NewUUID()
	
	_, err := DB.Exec(
		`INSERT INTO categorias (id, nombre, descripcion) VALUES (?, ?, ?)`,
		id, nombre, descripcion,
	)
	if err != nil {
		return nil, err
	}

	return GetCategoriaByID(id)
}

func GetCategoriaByID(id string) (*model.Categoria, error) {
	var c model.Categoria
	err := DB.QueryRow(
		`SELECT id, nombre, descripcion, created_at, updated_at 
		 FROM categorias WHERE id = ?`, id,
	).Scan(&c.ID, &c.Nombre, &c.Descripcion, &c.CreatedAt, &c.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetAllCategorias() ([]model.Categoria, error) {
	rows, err := DB.Query(
		`SELECT id, nombre, descripcion, created_at, updated_at FROM categorias ORDER BY nombre`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categorias []model.Categoria
	for rows.Next() {
		var c model.Categoria
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Descripcion, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}

func UpdateCategoria(id string, nombre string, descripcion *string) (*model.Categoria, error) {
	_, err := DB.Exec(
		`UPDATE categorias SET nombre = ?, descripcion = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		nombre, descripcion, id,
	)
	if err != nil {
		return nil, err
	}
	return GetCategoriaByID(id)
}

func DeleteCategoria(id string) error {
	_, err := DB.Exec(`DELETE FROM categorias WHERE id = ?`, id)
	return err
}

func CategoriaNameExists(nombre string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM categorias WHERE nombre = ?`, nombre).Scan(&count)
	return count > 0, err
}