package repository

import (
	"github.com/stockapp/backend/internal/model"
)

func CreateUsuario(email, passwordHash string, nombre *string, rol model.UserRole) (*model.Usuario, error) {
	// PostgreSQL generates UUID automatically via gen_random_uuid()
	var id string
	err := DB.QueryRow(
		`INSERT INTO usuarios (email, password_hash, nombre, rol) VALUES ($1, $2, $3, $4) RETURNING id`,
		email, passwordHash, nombre, string(rol),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetUsuarioByID(id)
}

func GetUsuarioByID(id string) (*model.Usuario, error) {
	var u model.Usuario
	err := DB.QueryRow(
		`SELECT id, email, password_hash, nombre, rol, created_at, updated_at 
		 FROM usuarios WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Nombre, &u.Rol, &u.CreatedAt, &u.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUsuarioByEmail(email string) (*model.Usuario, error) {
	var u model.Usuario
	err := DB.QueryRow(
		`SELECT id, email, password_hash, nombre, rol, created_at, updated_at 
		 FROM usuarios WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Nombre, &u.Rol, &u.CreatedAt, &u.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetAllUsuarios() ([]model.Usuario, error) {
	rows, err := DB.Query(
		`SELECT id, email, nombre, rol, created_at, updated_at FROM usuarios`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		if err := rows.Scan(&u.ID, &u.Email, &u.Nombre, &u.Rol, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

func UpdateUsuario(id string, nombre *string) (*model.Usuario, error) {
	_, err := DB.Exec(
		`UPDATE usuarios SET nombre = $1, updated_at = NOW() WHERE id = $2`,
		nombre, id,
	)
	if err != nil {
		return nil, err
	}
	return GetUsuarioByID(id)
}

func DeleteUsuario(id string) error {
	_, err := DB.Exec(`DELETE FROM usuarios WHERE id = $1`, id)
	return err
}

func EmailExists(email string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}