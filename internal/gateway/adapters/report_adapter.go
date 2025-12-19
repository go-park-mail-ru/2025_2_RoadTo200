package adapters

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

var _ service.SupportService = (*SupportServiceAdapter)(nil)

type SupportServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewSupportServiceAdapter(client pb.CoreServiceClient) service.SupportService {
	return &SupportServiceAdapter{
		client: client,
	}
}

// NewReportServiceAdapter is an alias for NewSupportServiceAdapter for consistency
func NewReportServiceAdapter(client pb.CoreServiceClient) service.SupportService {
	return NewSupportServiceAdapter(client)
}

func (a *SupportServiceAdapter) CreateTicket(userID uuid.UUID, request *dto.SupportTicketRequest) (*dto.SupportTicketResponse, error) {
	resp, err := a.client.CreateReport(context.Background(), &pb.CreateReportRequest{
		UserId:   userID.String(),
		Category: request.Category,
		Text:     request.Text,
		Email:    request.Email,
	})
	if err != nil {
		return nil, err
	}

	return protoToSupportTicketResponse(resp.Report), nil
}

func (a *SupportServiceAdapter) GetUserTickets(userID uuid.UUID, limit, offset int) (*dto.SupportTicketsListResponse, error) {
	resp, err := a.client.GetUserReports(context.Background(), &pb.GetUserReportsRequest{
		UserId: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	tickets := make([]dto.SupportTicketResponse, len(resp.Reports))
	for i, report := range resp.Reports {
		tickets[i] = *protoToSupportTicketResponse(report)
	}

	return &dto.SupportTicketsListResponse{
		Tickets: tickets,
		Total:   int(resp.Total),
	}, nil
}

func (a *SupportServiceAdapter) GetUserTicket(userID uuid.UUID, ticketID uuid.UUID) (*dto.SupportTicketDetailResponse, error) {
	resp, err := a.client.GetUserReport(context.Background(), &pb.GetUserReportRequest{
		UserId:   userID.String(),
		ReportId: ticketID.String(),
	})
	if err != nil {
		return nil, err
	}

	return protoToSupportTicketDetailResponseFromDetail(resp.Report), nil
}

func (a *SupportServiceAdapter) GetSupportStats() (*dto.SupportStatsResponse, error) {
	resp, err := a.client.GetSupportStats(context.Background(), &pb.GetSupportStatsRequest{})
	if err != nil {
		return nil, err
	}

	allTickets := make([]dto.SupportTicketDetailResponse, len(resp.AllTickets))
	for i, ticket := range resp.AllTickets {
		allTickets[i] = *protoToSupportTicketDetailResponseFromDetail(ticket)
	}

	return &dto.SupportStatsResponse{
		TotalTickets:        int(resp.TotalTickets),
		TicketsByCategory:   convertStringInt32Map(resp.TicketsByCategory),
		TicketsByStatus:     convertStringInt32Map(resp.TicketsByStatus),
		AverageResponseTime: resp.AverageResponseTime,
		AllTickets:          allTickets,
	}, nil
}

// Helper functions

func protoToSupportTicketResponse(pbReport *pb.Report) *dto.SupportTicketResponse {
	return &dto.SupportTicketResponse{
		ID:        pbReport.Id,
		Category:  pbReport.Category,
		Status:    pbReport.Status,
		CreatedAt: pbReport.CreatedAt.AsTime(),
	}
}

func protoToSupportTicketDetailResponse(pbReport *pb.Report) *dto.SupportTicketDetailResponse {
	return &dto.SupportTicketDetailResponse{
		ID:        pbReport.Id,
		Category:  pbReport.Category,
		Status:    pbReport.Status,
		CreatedAt: pbReport.CreatedAt.AsTime(),
		// Note: Text and Email are not in pb.Report, they should be in pb.ReportDetail
		// This will need to be fixed when proto is updated
	}
}

func protoToSupportTicketDetailResponseFromDetail(pbDetail *pb.ReportDetail) *dto.SupportTicketDetailResponse {
	return &dto.SupportTicketDetailResponse{
		ID:        pbDetail.Id,
		Category:  pbDetail.Category,
		Text:      pbDetail.Text,
		Email:     pbDetail.Email,
		Status:    pbDetail.Status,
		CreatedAt: pbDetail.CreatedAt.AsTime(),
	}
}

func convertStringInt32Map(m map[string]int32) map[string]int {
	result := make(map[string]int, len(m))
	for k, v := range m {
		result[k] = int(v)
	}
	return result
}
