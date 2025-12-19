package dto

//go:generate easyjson -all -no_std_marshalers auth_dto.go

// RegisterRequest represents registration request
// @Description Запрос для регистрации нового пользователя
type RegisterRequest struct {
	Email           string `json:"email" binding:"required" example:"user@example.com"`
	Password        string `json:"password" binding:"required" example:"securepassword123"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required" example:"securepassword123"`
}

// LoginRequest represents login request
// @Description Запрос для аутентификации пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"user@example.com"`
}

// LoginResponse represents login response
type LoginResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"user@example.com"`
}
