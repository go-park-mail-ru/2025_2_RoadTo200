package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type SupportService interface {
	// For users
	CreateTicket(userID uuid.UUID, request *dto.SupportTicketRequest) (*dto.SupportTicketResponse, error)
	GetUserTickets(userID uuid.UUID, limit, offset int) (*dto.SupportTicketsListResponse, error)
}
