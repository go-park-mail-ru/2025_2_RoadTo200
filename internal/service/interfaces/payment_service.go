package service

import (
	"context"

	"github.com/google/uuid"
)

// PaymentService defines payment service contract
type PaymentService interface {
	// CreatePaymentLink создает ссылку на оплату премиум-подписки
	CreatePaymentLink(ctx context.Context, userID uuid.UUID, planID string, amount float64) (string, error)

	// ProcessWebhook обрабатывает webhook от ЮMoney о поступлении платежа
	ProcessWebhook(ctx context.Context, webhookData map[string]string) error
}
