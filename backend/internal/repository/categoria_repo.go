package repository

import (
	"github.com/stockapp/backend/internal/model"
)

func CreateCategoria(nombre string, descripcion *string) (*model.Categoria, error) {
	// PostgreSQL generates UUID automatically via gen_random_uuid()
	var id string
	err := DB.QueryRow(
		`INSERT INTO categorias (nombre, descripcion) VALUES ($1, $2) RETURNING id`,
		nombre, descripcion,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetCategoriaByID(id)
}

func GetCategoriaByID(id string) (*model.Categoria, error) {
	var c model.Categoria
	err := DB.QueryRow(
		`SELECT id, nombre, descripcion, created_at, updated_at 
		 FROM categorias WHERE id = $1`, id,
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
		`UPDATE categorias SET nombre = $1, descripcion = $2, updated_at = NOW() WHERE id = $3`,
		nombre, descripcion, id,
	)
	if err != nil {
		return nil, err
	}
	return GetCategoriaByID(id)
}

func DeleteCategoria(id string) error {
	_, err := DB.Exec(`DELETE FROM categorias WHERE id = $1`, id)
	return err
}

func CategoriaNameExists(nombre string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM categorias WHERE nombre = $1`, nombre).Scan(&count)
	return count > 0, err
}