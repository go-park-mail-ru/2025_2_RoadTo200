package dto

//go:generate easyjson -all -no_std_marshalers notification_dto.go

import domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"

// NotificationsResponse ответ со списком уведомлений
type NotificationsResponse struct {
	Notifications []domain.Notification `json:"notifications"`
	Total         int                   `json:"total"`
	Limit         int                   `json:"limit"`
	Offset        int                   `json:"offset"`
}
