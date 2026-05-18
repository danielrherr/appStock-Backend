package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/stockapp/backend/internal/repository"
)

type NotificationService struct {
	firebaseAPIKey string
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		firebaseAPIKey: os.Getenv("FIREBASE_API_KEY"),
	}
}

type FCMMessage struct {
	Message FCMMessageContent `json:"message"`
}

type FCMMessageContent struct {
	Token       string             `json:"token,omitempty"`
	Topic       string             `json:"topic,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Notification FCMNotification   `json:"notification,omitempty"`
	Android      FCMAndroid        `json:"android,omitempty"`
	APNS         FCMAPNS           `json:"apns,omitempty"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type FCMAndroid struct {
	Priority string            `json:"priority"`
	Data     map[string]string `json:"data,omitempty"`
}

type FCMAPNS struct {
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type SendNotificationRequest struct {
	Token   string            `json:"token,omitempty"`
	Topic   string            `json:"topic,omitempty"`
	Title   string            `json:"title"`
	Body    string            `json:"body"`
	Data    map[string]string `json:"data,omitempty"`
}

// SendPushNotification envía una notificación a un token específico
func (s *NotificationService) SendPushNotification(req SendNotificationRequest) error {
	if s.firebaseAPIKey == "" {
		return fmt.Errorf("FIREBASE_API_KEY no configurada")
	}

	msg := FCMMessage{
		Message: FCMMessageContent{
			Token: req.Token,
			Notification: FCMNotification{
				Title: req.Title,
				Body:  req.Body,
			},
			Data: req.Data,
			Android: FCMAndroid{
				Priority: "high",
				Data:     req.Data,
			},
		},
	}

	return s.sendToFCM(msg)
}

// SendToTopic envía notificación a todos los dispositivos suscritos a un topic
func (s *NotificationService) SendToTopic(topic, title, body string, data map[string]string) error {
	if s.firebaseAPIKey == "" {
		return fmt.Errorf("FIREBASE_API_KEY no configurada")
	}

	msg := FCMMessage{
		Message: FCMMessageContent{
			Topic: topic,
			Notification: FCMNotification{
				Title: title,
				Body:  body,
			},
			Data: data,
		},
	}

	return s.sendToFCM(msg)
}

func (s *NotificationService) sendToFCM(msg FCMMessage) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"https://fcm.googleapis.com/v1/projects/YOUR_PROJECT_ID/messages:send",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.firebaseAPIKey)

	// Usar http.Client con timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM error: %d", resp.StatusCode)
	}

	return nil
}

// CheckStockAlerts verifica productos con stock bajo y envía notificaciones
func (s *NotificationService) CheckStockAlerts() error {
	productos, err := repository.GetProductosStockBajo()
	if err != nil {
		return err
	}

	if len(productos) == 0 {
		return nil
	}

	// Por cada producto con stock bajo, notificar a los usuarios
	for _, p := range productos {
		title := "⚠️ Alerta de Stock"
		body := fmt.Sprintf("El producto '%s' tiene stock bajo (%d)", p.Nombre, p.StockActual)
		data := map[string]string{
			"type":       "stock_alert",
			"product_id": p.ID,
		}

		// Enviar notificación de alerta (topic broadcasts)
		s.SendToTopic("stock-alerts", title, body, data)
	}

	return nil
}