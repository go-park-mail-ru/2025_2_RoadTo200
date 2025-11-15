package dto

import (
	"time"
)

// SupportTicketRequest represents support ticket creation request
type SupportTicketRequest struct {
	Category string `json:"category" example:"technical" binding:"required" enums:"technical,feature,question,security,billing,device"`
	Text     string `json:"text" example:"При попытке загрузить фото приложение вылетает" binding:"required"`
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
}

type SupportTicketResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Category  string    `json:"category" example:"technical"`
	Status    string    `json:"status" example:"open" enums:"open,work,closed"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

// SupportTicketsListResponse represents list of support tickets
type SupportTicketsListResponse struct {
	Tickets []SupportTicketResponse `json:"tickets"`
	Total   int                     `json:"total" example:"5"`
}

// SupportStatsResponse represents support statistics
type SupportStatsResponse struct {
	TotalTickets        int            `json:"total_tickets" example:"150"`
	TicketsByCategory   map[string]int `json:"tickets_by_category" example:"technical:50,feature:40,question:35,security:15,billing:10,device:5"`
	TicketsByStatus     map[string]int `json:"tickets_by_status" example:"open:25,in_progress:10,closed:115"`
	AverageResponseTime string         `json:"average_response_time" example:"2h30m"`
}
