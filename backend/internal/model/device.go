package model

import (
	"time"
)

type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformWeb     Platform = "web"
)

type Device struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id"`
	Token     string    `json:"token"`
	Platform  Platform  `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterDeviceRequest struct {
	Token    string  `json:"token" validate:"required"`
	Platform Platform `json:"platform" validate:"required,oneof=android ios web"`
}