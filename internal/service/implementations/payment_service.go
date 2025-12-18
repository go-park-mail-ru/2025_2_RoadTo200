package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

type PaymentService struct {
	userRepo         interfaces.UserRepository
	subscriptionRepo interfaces.SubscriptionRepository
	config           *config.Config
	logger           logger.Log
}

func NewPaymentService(
	userRepo interfaces.UserRepository,
	subscriptionRepo interfaces.SubscriptionRepository,
	cfg *config.Config,
	l logger.Log,
) *PaymentService {
	return &PaymentService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		config:           cfg,
		logger:           l,
	}
}

// CreatePaymentLink создает ссылку на оплату премиум-подписки
func (s *PaymentService) CreatePaymentLink(ctx context.Context, userID uuid.UUID, planID string, amount float64) (string, error) {
	s.logger.Tracef("CreatePaymentLink: userID=%s, planID=%s, amount=%.2f", userID, planID, amount)

	// Валидация плана
	planType := constants.PlanType(planID)
	if planType != constants.PlanTypeWeek && planType != constants.PlanTypeMonth && planType != constants.PlanTypeQuarter {
		return "", fmt.Errorf("invalid plan ID: %s", planID)
	}

	// Проверяем, что у пользователя нет активной подписки (фронт должен это блокировать, но проверим на всякий случай)
	activeSub, err := s.subscriptionRepo.GetActiveSubscription(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to check active subscription: %w", err)
	}
	if activeSub != nil {
		return "", fmt.Errorf("user already has active subscription")
	}

	// Хардкод суммы для тестов - всегда 1 рубль
	testAmount := 2.00

	// Создаем label для идентификации платежа: userID:planID:timestamp
	timestamp := time.Now().Unix()
	label := fmt.Sprintf("%s:%s:%d", userID.String(), planID, timestamp)

	// Формируем URL для quickpay
	baseURL := "https://yoomoney.ru/quickpay/confirm"
	params := url.Values{}
	params.Set("receiver", s.config.Payment.Receiver)
	params.Set("quickpay-form", "button")
	params.Set("paymentType", "AC") // AC - с банковской карты, PC - из кошелька. Можем сделать опциональным
	params.Set("sum", fmt.Sprintf("%.2f", testAmount))
	params.Set("label", label)

	paymentURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	s.logger.Debugf("Created payment link for user %s: %s", userID, paymentURL)
	return paymentURL, nil
}

// ProcessWebhook обрабатывает webhook от ЮMoney о поступлении платежа
func (s *PaymentService) ProcessWebhook(ctx context.Context, webhookData map[string]string) error {
	s.logger.Infof("ProcessWebhook: processing webhook from YooMoney")
	s.logger.Infof("ProcessWebhook: webhook data: %+v", webhookData)

	// Проверяем подпись
	if !s.verifyWebhookSignature(webhookData) {
		s.logger.Error("ProcessWebhook: invalid signature - webhook rejected")
		return fmt.Errorf("invalid webhook signature")
	}

	// Извлекаем данные из webhook
	notificationType := webhookData["notification_type"]
	if notificationType != "p2p-incoming" && notificationType != "card-incoming" {
		s.logger.Warnf("ProcessWebhook: unsupported notification type: %s", notificationType)
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	label := webhookData["label"]
	if label == "" {
		s.logger.Warn("ProcessWebhook: empty label")
		return fmt.Errorf("empty label in webhook")
	}

	// Парсим label: userID:planID:timestamp
	parts := strings.Split(label, ":")
	if len(parts) != 3 {
		s.logger.Warnf("ProcessWebhook: invalid label format: %s", label)
		return fmt.Errorf("invalid label format")
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		s.logger.Warnf("ProcessWebhook: invalid userID in label: %s", parts[0])
		return fmt.Errorf("invalid userID in label: %w", err)
	}

	planID := parts[1]
	planType := constants.PlanType(planID)
	if planType != constants.PlanTypeWeek && planType != constants.PlanTypeMonth && planType != constants.PlanTypeQuarter {
		s.logger.Warnf("ProcessWebhook: invalid planID in label: %s", planID)
		return fmt.Errorf("invalid planID in label: %s", planID)
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found: %s", userID)
	}

	// Вычисляем даты подписки
	startDate := time.Now()
	var endDate time.Time
	switch planType {
	case constants.PlanTypeWeek:
		endDate = startDate.AddDate(0, 0, 7)
	case constants.PlanTypeMonth:
		endDate = startDate.AddDate(0, 1, 0)
	case constants.PlanTypeQuarter:
		endDate = startDate.AddDate(0, 3, 0)
	}

	// Обновляем is_premium и добавляем 6 суперлайков при активации премиума
	user.IsPremium = true
	// Добавляем 6 суперлайков к текущему количеству (не заменяем)
	user.SuperLikesCount += 6
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user premium status: %w", err)
	}
	s.logger.Infof("Updated user premium status: userID=%s, is_premium=true, super_likes_count=%d", userID, user.SuperLikesCount)

	// Создаем или обновляем подписку
	subscription := &domain.Subscription{
		UserID:    userID,
		PlanType:  planType,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	existingSub, err := s.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing subscription: %w", err)
	}

	if existingSub != nil {
		// Обновляем существующую подписку
		err = s.subscriptionRepo.Update(ctx, subscription)
		if err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
		s.logger.Infof("Updated subscription for user %s: plan=%s, end_date=%s", userID, planID, endDate)
	} else {
		// Создаем новую подписку
		err = s.subscriptionRepo.Create(ctx, subscription)
		if err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
		s.logger.Infof("Created subscription for user %s: plan=%s, end_date=%s", userID, planID, endDate)
	}

	return nil
}

// verifyWebhookSignature проверяет подпись webhook от ЮMoney
func (s *PaymentService) verifyWebhookSignature(webhookData map[string]string) bool {
	secret := s.config.Payment.NotificationSecret
	if secret == "" {
		s.logger.Error("Notification secret is not configured (MONEY_NOTIFICATION_SECRET env variable is empty)")
		return false
	}

	// Формируем строку для проверки подписи согласно документации:
	// notification_type&operation_id&amount&currency&datetime&sender&codepro&notification_secret&label
	notificationType := webhookData["notification_type"]
	operationID := webhookData["operation_id"]
	amount := webhookData["amount"]
	currency := webhookData["currency"]
	datetime := webhookData["datetime"]
	sender := webhookData["sender"]
	codepro := webhookData["codepro"]
	label := webhookData["label"]

	// Если label пустой, используем пустую строку
	if label == "" {
		label = ""
	}

	signatureString := fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s&%s&%s",
		notificationType, operationID, amount, currency, datetime, sender, codepro, secret, label)

	// Вычисляем SHA-1 хэш
	hash := sha1.Sum([]byte(signatureString))
	calculatedHash := hex.EncodeToString(hash[:])

	// Сравниваем с переданным хэшем
	receivedHash := webhookData["sha1_hash"]

	s.logger.Infof("Webhook signature verification:")
	s.logger.Infof("  Signature string: %s", signatureString)
	s.logger.Infof("  Calculated hash: %s", calculatedHash)
	s.logger.Infof("  Received hash: %s", receivedHash)

	isValid := calculatedHash == receivedHash
	if !isValid {
		s.logger.Warnf("Webhook signature verification failed!")
	} else {
		s.logger.Infof("Webhook signature verification successful")
	}

	return isValid
}
