package adapters

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/google/uuid"
)

var _ service.AuthService = (*AuthServiceAdapter)(nil)

type AuthServiceAdapter struct {
	client auth.AuthServiceClient
}

func NewAuthServiceAdapter(client auth.AuthServiceClient) service.AuthService {
	return &AuthServiceAdapter{
		client: client,
	}
}

func (a *AuthServiceAdapter) Register(email, password, passwordConfirm string) (*domain.User, *domain.Session, error) {
	resp, err := a.client.Register(context.Background(), &auth.RegisterRequest{
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	})
	if err != nil {
		return nil, nil, err
	}

	return protoToUser(resp.User), protoToSession(resp.Session), nil
}

func (a *AuthServiceAdapter) Login(email, password string) (*domain.User, *domain.Session, error) {
	resp, err := a.client.Login(context.Background(), &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, nil, err
	}

	return protoToUser(resp.User), protoToSession(resp.Session), nil
}

func (a *AuthServiceAdapter) Logout(token string) error {
	_, err := a.client.Logout(context.Background(), &auth.LogoutRequest{
		Token: token,
	})
	return err
}

func (a *AuthServiceAdapter) ValidateSession(token string) (*domain.User, error) {
	resp, err := a.client.ValidateSession(context.Background(), &auth.ValidateSessionRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return protoToUser(resp.User), nil
}

func (a *AuthServiceAdapter) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	resp, err := a.client.GetUserByID(context.Background(), &auth.GetUserByIDRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return protoToUser(resp.User), nil
}

// Helper functions to convert Proto -> Domain

func protoToUser(pbUser *auth.User) *domain.User {
	if pbUser == nil {
		return nil
	}

	id, _ := uuid.Parse(pbUser.Id)

	user := &domain.User{
		ID:         id,
		Email:      pbUser.Email,
		Name:       pbUser.Name,
		BirthDate:  pbUser.BirthDate.AsTime(),
		IsVerified: pbUser.IsVerified,
		LastActive: pbUser.LastActive.AsTime(),
		CreatedAt:  pbUser.CreatedAt.AsTime(),
		UpdatedAt:  pbUser.UpdatedAt.AsTime(),
	}

	// Handle optional fields
	if pbUser.Phone != nil {
		user.Phone = pbUser.Phone
	}
	if pbUser.Bio != nil {
		user.Bio = pbUser.Bio
	}
	if pbUser.City != nil {
		user.City = pbUser.City
	}
	if pbUser.Artist != nil {
		user.Artist = pbUser.Artist
	}
	if pbUser.Quote != nil {
		user.Quote = pbUser.Quote
	}

	user.Gender = constants.Gender(pbUser.Gender)

	// Gender conversion (simplified for now, assuming string match or mapping needed if enum)
	// Assuming domain has Gender as string or compatible type.
	// If domain.Gender is custom type, we might need casting.
	// user.Gender = constants.Gender(pbUser.Gender)
	// Let's assume for now it's compatible or we need to check domain definition.
	// Checking domain definition from previous context: Gender is constants.Gender (string alias likely)

	return user
}

func protoToSession(pbSession *auth.Session) *domain.Session {
	if pbSession == nil {
		return nil
	}

	return &domain.Session{
		Token:     pbSession.Token,
		UserEmail: pbSession.UserEmail,
		ExpiresAt: pbSession.ExpiresAt.AsTime(),
	}
}
