package service

import (
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/repository"
)

type DeviceService struct{}

func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

func (s *DeviceService) RegisterDevice(userID *string, token string, platform model.Platform) (*model.Device, error) {
	// Verificar si el token ya existe
	exists, err := repository.DeviceTokenExists(token)
	if err != nil {
		return nil, err
	}

	if exists {
		// Token existente - actualizar user_id si es diferente
		devices, err := repository.GetAllDevices()
		if err != nil {
			return nil, err
		}
		// Buscar y actualizar el device con este token
		for _, d := range devices {
			if d.Token == token {
				return repository.UpdateDeviceToken(d.ID, token)
			}
		}
	}

	return repository.CreateDevice(userID, token, platform)
}

func (s *DeviceService) GetUserDevices(userID string) ([]model.Device, error) {
	return repository.GetDevicesByUser(userID)
}

func (s *DeviceService) DeleteDevice(id string) error {
	return repository.DeleteDevice(id)
}