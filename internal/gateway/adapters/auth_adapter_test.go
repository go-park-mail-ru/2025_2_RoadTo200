package adapters

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthServiceAdapter_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	email := "test@example.com"
	password := "password"
	passwordConfirm := "password"
	userID := uuid.New().String()
	token := "test_token"

	resp := &auth.RegisterResponse{
		User: &auth.User{
			Id:         userID,
			Email:      email,
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			BirthDate:  timestamppb.New(time.Now().AddDate(-20, 0, 0)),
			LastActive: timestamppb.Now(),
		},
		Session: &auth.Session{
			Token:     token,
			UserEmail: email,
			ExpiresAt: timestamppb.Now(),
		},
	}

	mockClient.EXPECT().
		Register(gomock.Any(), &auth.RegisterRequest{
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
		}).
		Return(resp, nil)

	user, session, err := adapter.Register(context.Background(), email, password, passwordConfirm)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, token, session.Token)
}

func TestAuthServiceAdapter_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	mockClient.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("grpc error"))

	user, session, err := adapter.Register(context.Background(), "email", "pass", "pass")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthServiceAdapter_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	email := "test@example.com"
	password := "password"
	userID := uuid.New().String()
	token := "test_token"

	resp := &auth.LoginResponse{
		User: &auth.User{
			Id:         userID,
			Email:      email,
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			BirthDate:  timestamppb.New(time.Now().AddDate(-20, 0, 0)),
			LastActive: timestamppb.Now(),
		},
		Session: &auth.Session{
			Token:     token,
			UserEmail: email,
			ExpiresAt: timestamppb.Now(),
		},
	}

	mockClient.EXPECT().
		Login(gomock.Any(), &auth.LoginRequest{
			Email:    email,
			Password: password,
		}).
		Return(resp, nil)

	user, session, err := adapter.Login(context.Background(), email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, token, session.Token)
}

func TestAuthServiceAdapter_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	token := "test_token"

	mockClient.EXPECT().
		Logout(gomock.Any(), &auth.LogoutRequest{
			Token: token,
		}).
		Return(&auth.LogoutResponse{Success: true}, nil)

	err := adapter.Logout(context.Background(), token)
	assert.NoError(t, err)
}

func TestAuthServiceAdapter_ValidateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	token := "test_token"
	userID := uuid.New().String()

	resp := &auth.ValidateSessionResponse{
		User: &auth.User{
			Id:         userID,
			Email:      "test@example.com",
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			BirthDate:  timestamppb.New(time.Now().AddDate(-20, 0, 0)),
			LastActive: timestamppb.Now(),
		},
	}

	mockClient.EXPECT().
		ValidateSession(gomock.Any(), &auth.ValidateSessionRequest{
			Token: token,
		}).
		Return(resp, nil)

	user, err := adapter.ValidateSession(context.Background(), token)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID.String())
}

func TestAuthServiceAdapter_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	adapter := NewAuthServiceAdapter(mockClient)

	userID := uuid.New()

	resp := &auth.GetUserByIDResponse{
		User: &auth.User{
			Id:         userID.String(),
			Email:      "test@example.com",
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			BirthDate:  timestamppb.New(time.Now().AddDate(-20, 0, 0)),
			LastActive: timestamppb.Now(),
		},
	}

	mockClient.EXPECT().
		GetUserByID(gomock.Any(), &auth.GetUserByIDRequest{
			UserId: userID.String(),
		}).
		Return(resp, nil)

	user, err := adapter.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
}
