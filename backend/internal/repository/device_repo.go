package repository

import (
	"github.com/stockapp/backend/internal/model"
)

func CreateDevice(userID *string, token string, platform model.Platform) (*model.Device, error) {
	// PostgreSQL generates UUID automatically via gen_random_uuid()
	var id string
	err := DB.QueryRow(
		`INSERT INTO devices (user_id, token, platform) VALUES ($1, $2, $3) RETURNING id`,
		userID, token, string(platform),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetDeviceByID(id)
}

func GetDeviceByID(id string) (*model.Device, error) {
	var d model.Device
	err := DB.QueryRow(
		`SELECT id, user_id, token, platform, created_at FROM devices WHERE id = $1`, id,
	).Scan(&d.ID, &d.UserID, &d.Token, &d.Platform, &d.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func GetDevicesByUser(userID string) ([]model.Device, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, token, platform, created_at FROM devices WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Token, &d.Platform, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func GetAllDevices() ([]model.Device, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, token, platform, created_at FROM devices`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Token, &d.Platform, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func UpdateDeviceToken(id, token string) (*model.Device, error) {
	_, err := DB.Exec(
		`UPDATE devices SET token = $1 WHERE id = $2`,
		token, id,
	)
	if err != nil {
		return nil, err
	}
	return GetDeviceByID(id)
}

func DeleteDevice(id string) error {
	_, err := DB.Exec(`DELETE FROM devices WHERE id = $1`, id)
	return err
}

func DeviceTokenExists(token string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM devices WHERE token = $1`, token).Scan(&count)
	return count > 0, err
}