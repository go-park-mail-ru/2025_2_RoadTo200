package service

import (
	"fmt"
	//"log"
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
	reportRepo      interfaces.ReportRepository
	logger          logger.Log
	telegramService TelegramService
}

func NewSupportService(
	reportRepo interfaces.ReportRepository,
	l logger.Log,
	telegramService TelegramService,
) service.SupportService {
	return &supportService{
		reportRepo:      reportRepo,
		logger:          l,
		telegramService: telegramService,
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
		Comment:   nil,
		Status:    constants.DefaultReportStatus(),
		CreatedAt: time.Now(),
		WorkAt:    nil,
		ClosedAt:  nil,
	}

	// Сохраняем в БД
	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}

	go func() {
		messageID, err := s.telegramService.SendNewTicketNotification(report)
		if err != nil {
			s.logger.Errorf("Failed to send telegram notification: %v", err)
		} else {
			// Сохраняем ID сообщения если нужно (опционально)
			s.logger.Infof("Telegram notification sent with message ID: %d", messageID)
		}
	}()

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

	// Получаем все обращения для детальной статистики
	allReports, err := s.reportRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all reports: %w", err)
	}

	// Конвертируем в DTO
	stats := s.convertStatsToDTO(dbStats, allReports)

	return stats, nil
}

// convertStatsToDTO конвертирует статистику из БД в DTO
func (s *supportService) convertStatsToDTO(dbStats *interfaces.ReportStatistics, allReports []domain.Report) *dto.SupportStatsResponse {
	// Конвертируем все обращения в DTO
	allTickets := make([]dto.SupportTicketDetailResponse, 0, len(allReports))
	for _, report := range allReports {
		allTickets = append(allTickets, *s.convertToDetailResponse(&report))
	}

	return &dto.SupportStatsResponse{
		TotalTickets: int(dbStats.TotalTickets),
		TicketsByCategory: map[string]int{
			"technical": int(dbStats.TicketsByTheme.Technical),
			"feature":   int(dbStats.TicketsByTheme.Feature),
			"question":  int(dbStats.TicketsByTheme.Question),
			"security":  int(dbStats.TicketsByTheme.Security),
			"billing":   int(dbStats.TicketsByTheme.Billing),
			"device":    int(dbStats.TicketsByTheme.Device),
		},
		TicketsByStatus: map[string]int{
			"open":        int(dbStats.TicketsByStatus.Open),
			"in_progress": int(dbStats.TicketsByStatus.Work),
			"closed":      int(dbStats.TicketsByStatus.Closed),
		},
		AverageResponseTime: dbStats.AverageResponseTime,
		AllTickets:          allTickets,
	}
}

// GetUserTicket возвращает детальную информацию об обращении
func (s *supportService) GetUserTicket(userID uuid.UUID, ticketID uuid.UUID) (*dto.SupportTicketDetailResponse, error) {
	// Получаем обращение из БД
	report, err := s.reportRepo.GetById(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Проверяем, что обращение принадлежит пользователю
	if report.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Конвертируем в детальный ответ
	return s.convertToDetailResponse(report), nil
}

// convertToDetailResponse конвертирует доменную сущность в детальный ответ
func (s *supportService) convertToDetailResponse(report *domain.Report) *dto.SupportTicketDetailResponse {
	return &dto.SupportTicketDetailResponse{
		ID:        report.ID.String(),
		Category:  string(report.Theme),
		Text:      report.Problem,
		Email:     report.Contact,
		Status:    string(report.Status),
		CreatedAt: report.CreatedAt,
	}
}
