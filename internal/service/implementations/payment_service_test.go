package service

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentService_CreatePaymentLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	subscriptionRepo := mocks.NewMockSubscriptionRepository(ctrl)
	logger := mocks.NewMockLogger()

	cfg := &config.Config{
		Payment: config.PaymentConfig{
			Receiver: "410011234567890",
		},
	}

	service := NewPaymentService(userRepo, subscriptionRepo, cfg, logger)

	userID := uuid.New()

	tests := []struct {
		name          string
		planID        string
		amount        float64
		setupMocks    func()
		expectedError bool
		expectedURL   bool
	}{
		{
			name:   "successful week plan",
			planID: "week",
			amount: 100.0,
			setupMocks: func() {
				subscriptionRepo.EXPECT().
					GetActiveSubscription(gomock.Any(), userID).
					Return(nil, nil) // No active subscription
			},
			expectedError: false,
			expectedURL:   true,
		},
		{
			name:   "successful month plan",
			planID: "month",
			amount: 200.0,
			setupMocks: func() {
				subscriptionRepo.EXPECT().
					GetActiveSubscription(gomock.Any(), userID).
					Return(nil, nil) // No active subscription
			},
			expectedError: false,
			expectedURL:   true,
		},
		{
			name:   "successful quarter plan",
			planID: "quarter",
			amount: 300.0,
			setupMocks: func() {
				subscriptionRepo.EXPECT().
					GetActiveSubscription(gomock.Any(), userID).
					Return(nil, nil) // No active subscription
			},
			expectedError: false,
			expectedURL:   true,
		},
		{
			name:   "invalid plan ID",
			planID: "invalid",
			amount: 100.0,
			setupMocks: func() {
			},
			expectedError: true,
			expectedURL:   false,
		},
		{
			name:   "user has active subscription",
			planID: "month",
			amount: 200.0,
			setupMocks: func() {
				subscriptionRepo.EXPECT().
					GetActiveSubscription(gomock.Any(), userID).
					Return(&domain.Subscription{}, nil) // Has active subscription
			},
			expectedError: true,
			expectedURL:   false,
		},
		{
			name:   "subscription check error",
			planID: "month",
			amount: 200.0,
			setupMocks: func() {
				subscriptionRepo.EXPECT().
					GetActiveSubscription(gomock.Any(), userID).
					Return(nil, assert.AnError)
			},
			expectedError: true,
			expectedURL:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url, err := service.CreatePaymentLink(context.Background(), userID, tt.planID, tt.amount)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				if tt.expectedURL {
					assert.NotEmpty(t, url)
					assert.Contains(t, url, "https://yoomoney.ru/quickpay/confirm")
					assert.Contains(t, url, "receiver=410011234567890")
					assert.Contains(t, url, "sum=2.00") // Hardcoded test amount
					assert.Contains(t, url, userID.String())
					assert.Contains(t, url, tt.planID)
				}
			}
		})
	}
}

func TestPaymentService_ProcessWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	subscriptionRepo := mocks.NewMockSubscriptionRepository(ctrl)
	logger := mocks.NewMockLogger()

	cfg := &config.Config{
		Payment: config.PaymentConfig{
			Receiver:           "410011234567890",
			NotificationSecret: "test_secret",
		},
	}

	service := NewPaymentService(userRepo, subscriptionRepo, cfg, logger)

	userID := uuid.MustParse("8a721410-1e91-4531-b2c6-ed5a65c7fd47") // Fixed UUID for deterministic tests

	tests := []struct {
		name          string
		webhookData   map[string]string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful webhook processing - week plan",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:              userID,
						IsPremium:       false,
						SuperLikesCount: 0,
					}, nil)
				userRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				subscriptionRepo.EXPECT().
					GetByUserID(gomock.Any(), userID).
					Return(nil, nil) // No existing subscription
				subscriptionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "successful webhook processing - month plan with existing subscription",
			webhookData: map[string]string{
				"notification_type": "card-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":month:1234567890",
				"sha1_hash":         "672a9410e6651c8e5649e1d0325c3c051466105b",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:              userID,
						IsPremium:       false,
						SuperLikesCount: 0,
					}, nil)
				userRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				subscriptionRepo.EXPECT().
					GetByUserID(gomock.Any(), userID).
					Return(&domain.Subscription{}, nil) // Has existing subscription
				subscriptionRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "invalid webhook signature",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "invalid_hash",
			},
			setupMocks: func() {
			},
			expectedError: true,
		},
		{
			name: "unsupported notification type",
			webhookData: map[string]string{
				"notification_type": "unsupported",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "f271869803f9ee43a3ab340b23492c3602ff41b0",
			},
			setupMocks: func() {
			},
			expectedError: true,
		},
		{
			name: "empty label",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             "",
				"sha1_hash":         "beefe1ac62d9c150eede3dd451a382e7ebf3139e",
			},
			setupMocks: func() {
			},
			expectedError: true,
		},
		{
			name: "invalid label format",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             "invalid-label-format",
				"sha1_hash":         "274b85ff97347431231893f60bc31dd4cf614ecd",
			},
			setupMocks: func() {
			},
			expectedError: true,
		},
		{
			name: "invalid user ID in label",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             "invalid-uuid:week:1234567890",
				"sha1_hash":         "30235039ea10b8cc51a96aee13d2c5232d9afb3d",
			},
			setupMocks: func() {
			},
			expectedError: true,
		},
		{
			name: "user not found",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(nil, nil) // User not found
			},
			expectedError: true,
		},
		{
			name: "user repository error",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(nil, assert.AnError)
			},
			expectedError: true,
		},
		{
			name: "user update error",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:              userID,
						IsPremium:       false,
						SuperLikesCount: 0,
					}, nil)
				userRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			expectedError: true,
		},
		{
			name: "subscription repository error",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:              userID,
						IsPremium:       false,
						SuperLikesCount: 0,
					}, nil)
				userRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				subscriptionRepo.EXPECT().
					GetByUserID(gomock.Any(), userID).
					Return(nil, assert.AnError)
			},
			expectedError: true,
		},
		{
			name: "subscription create error",
			webhookData: map[string]string{
				"notification_type": "p2p-incoming",
				"operation_id":      "123456789",
				"amount":            "2.00",
				"currency":          "RUB",
				"datetime":          "2024-01-01T12:00:00Z",
				"sender":            "sender@example.com",
				"codepro":           "false",
				"label":             userID.String() + ":week:1234567890",
				"sha1_hash":         "311d3375e0bcb964a4d7e93903805aec0472dc8c",
			},
			setupMocks: func() {
				userRepo.EXPECT().
					GetByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:              userID,
						IsPremium:       false,
						SuperLikesCount: 0,
					}, nil)
				userRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
				subscriptionRepo.EXPECT().
					GetByUserID(gomock.Any(), userID).
					Return(nil, nil) // No existing subscription
				subscriptionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := service.ProcessWebhook(context.Background(), tt.webhookData)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
