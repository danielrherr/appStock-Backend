package repository

import (
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/utils"
)

func CreateUsuario(email, passwordHash string, nombre *string, rol model.UserRole) (*model.Usuario, error) {
	id := utils.NewUUID()
	
	_, err := DB.Exec(
		`INSERT INTO usuarios (id, email, password_hash, nombre, rol) VALUES (?, ?, ?, ?, ?)`,
		id, email, passwordHash, nombre, string(rol),
	)
	if err != nil {
		return nil, err
	}

	return GetUsuarioByID(id)
}

func GetUsuarioByID(id string) (*model.Usuario, error) {
	var u model.Usuario
	err := DB.QueryRow(
		`SELECT id, email, password_hash, nombre, rol, created_at, updated_at 
		 FROM usuarios WHERE id = ?`, id,
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
		 FROM usuarios WHERE email = ?`, email,
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
		`UPDATE usuarios SET nombre = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		nombre, id,
	)
	if err != nil {
		return nil, err
	}
	return GetUsuarioByID(id)
}

func DeleteUsuario(id string) error {
	_, err := DB.Exec(`DELETE FROM usuarios WHERE id = ?`, id)
	return err
}

func EmailExists(email string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE email = ?`, email).Scan(&count)
	return count > 0, err
}