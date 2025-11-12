package dto

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
