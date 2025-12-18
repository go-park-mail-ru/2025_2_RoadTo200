package dto

//go:generate easyjson -all -no_std_marshalers session_dto.go

// SessionResponse represents session check response
type SessionResponse struct {
	Authenticated bool                 `json:"authenticated" example:"true"`
	User          *SessionUserResponse `json:"user,omitempty"`
}

// SessionUserResponse represents user data in session response
type SessionUserResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"Алексей"`
}
