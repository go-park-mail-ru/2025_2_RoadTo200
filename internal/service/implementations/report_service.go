package service

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type supportService struct {
	reportRepo interfaces.ReportRepository
	logger     logger.Log
}

func NewSupportService(
	reportRepo interfaces.ReportRepository,
	l logger.Log,
) service.SupportService {
	return &supportService{
		reportRepo: reportRepo,
		logger:     l,
	}
}

// CreateTicket создает новое обращение в поддержку
func (s *supportService) CreateTicket(userID uuid.UUID, request *dto.SupportTicketRequest) (*dto.SupportTicketResponse, error) {
	// Конвертируем человекочитаемую категорию в техническую
	theme, err := constants.HumanCategoryToTheme(request.Category)
	if err != nil {
		return nil, err
	}

	// Создаем доменную сущность
	report := &domain.Report{
		ID:        uuid.New(),
		UserID:    userID,
		Theme:     theme,
		Problem:   request.Text,
		Contact:   request.Email,
		Comment:   "",
		Status:    constants.DefaultReportStatus(),
		CreatedAt: time.Now(),
		WorkAt:    time.Time{},
		ClosedAt:  time.Time{},
	}

	// Сохраняем в БД
	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}

	s.logger.Infof("Created report with ID: %v for userID: %v", report.ID, userID)
	return s.convertToResponse(report), nil
}

// GetUserTickets возвращает обращения пользователя
func (s *supportService) GetUserTickets(userID uuid.UUID, limit, offset int) (*dto.SupportTicketsListResponse, error) {
	reports, err := s.reportRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	// Применяем пагинацию
	start := offset
	if start > len(reports) {
		start = len(reports)
	}
	end := start + limit
	if end > len(reports) {
		end = len(reports)
	}

	paginatedReports := reports[start:end]

	// Конвертируем в ответ
	tickets := make([]dto.SupportTicketResponse, 0, len(paginatedReports))
	for _, report := range paginatedReports {
		tickets = append(tickets, *s.convertToResponse(&report))
	}

	s.logger.Infof("Fetched %d reports for userID: %v", len(tickets), userID)
	return &dto.SupportTicketsListResponse{
		Tickets: tickets,
		Total:   len(reports),
	}, nil
}

// convertToResponse конвертирует доменную сущность в ответ с человекочитаемыми названиями
func (s *supportService) convertToResponse(report *domain.Report) *dto.SupportTicketResponse {
	return &dto.SupportTicketResponse{
		ID:        report.ID.String(),
		Category:  constants.ThemeToHumanCategory(report.Theme),        // Возвращаем человекочитаемое название
		Status:    string(constants.StatusDisplayNames[report.Status]), // Возвращаем человекочитаемый статус
		CreatedAt: report.CreatedAt,
	}
}

// GetSupportStats возвращает статистику по всем обращениям
func (s *supportService) GetSupportStats() (*dto.SupportStatsResponse, error) {
	// Получаем статистику из репозитория
	dbStats, err := s.reportRepo.GetStatistics()
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	// Конвертируем в DTO
	stats := s.convertStatsToDTO(dbStats)

	return stats, nil
}

// convertStatsToDTO конвертирует статистику из БД в DTO
func (s *supportService) convertStatsToDTO(dbStats *interfaces.ReportStatistics) *SupportStatsResponse {
	return &dto.SupportStatsResponse{
		TotalTickets: dbStats.TotalTickets,
		TicketsByCategory: map[string]int{
			"technical": dbStats.TicketsByTheme.Technical,
			"feature":   dbStats.TicketsByTheme.Feature,
			"question":  dbStats.TicketsByTheme.Question,
			"security":  dbStats.TicketsByTheme.Security,
			"billing":   dbStats.TicketsByTheme.Billing,
			"device":    dbStats.TicketsByTheme.Device,
		},
		TicketsByStatus: map[string]int{
			"open":        dbStats.TicketsByStatus.Active,
			"in_progress": dbStats.TicketsByStatus.Work,
			"closed":      dbStats.TicketsByStatus.Closed,
		},
		AverageResponseTime: dbStats.AverageResponseTime,
	}
}
