package dto

//go:generate easyjson -all -no_std_marshalers base_dto.go

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
