package server

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/converters"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	logger      logger.Log
}

func NewAuthServer(authService service.AuthService, logger logger.Log) *AuthServer {
	return &AuthServer{
		authService: authService,
		logger:      logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.logger.Info("gRPC: Register request")

	user, session, err := s.authService.Register(ctx, req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		s.logger.Error("Register error: " + err.Error())
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}

	return &pb.RegisterResponse{
		User:    converters.UserToProto(user),
		Session: converters.SessionToProto(session),
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("gRPC: Login request")

	user, session, err := s.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Error("Login error: " + err.Error())
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	return &pb.LoginResponse{
		User:    converters.UserToProto(user),
		Session: converters.SessionToProto(session),
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	s.logger.Info("gRPC: Logout request")

	err := s.authService.Logout(ctx, req.Token)
	if err != nil {
		s.logger.Error("Logout error: " + err.Error())
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *AuthServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	s.logger.Debug("gRPC: ValidateSession request")

	user, err := s.authService.ValidateSession(ctx, req.Token)
	if err != nil {
		s.logger.Error("ValidateSession error: " + err.Error())
		return &pb.ValidateSessionResponse{
			User:    nil,
			IsValid: false,
		}, nil
	}

	return &pb.ValidateSessionResponse{
		User:    converters.UserToProto(user),
		IsValid: true,
	}, nil
}

func (s *AuthServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	s.logger.Debug("gRPC: GetUserByID request")

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := s.authService.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("GetUserByID error: " + err.Error())
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.GetUserByIDResponse{
		User: converters.UserToProto(user),
	}, nil
}
