package domain

import "github.com/google/uuid"

type Screen struct {
	ID       string    `json:"id"`
	ReportId uuid.UUID `json:"report_id"`
	Url      string    `json:"url"`
}
