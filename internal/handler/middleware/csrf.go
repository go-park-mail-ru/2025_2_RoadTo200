package middleware

import (
	"fmt"
	"net/http"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"google.golang.org/grpc/credentials/jwt"
)

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (*JwtToken, error) {
	return &JwtToken{Secret: []byte(secret)}, nil
}

type JwtCsrfClaims struct {
	Token     string
	UserEmail string
	ExpiresAt time.Time
}

func (tk *JwtToken) Create(s *domain.Session, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		Token:     s.Token,
		UserEmail: s.UserEmail,
		ExpiresAt: time.Now().Add(time.Duration(tokenExpTime) * time.Second),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(s *Session, inputToken string) (bool, error) {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.parseSecretGetter)
	if err != nil {
		return false, fmt.Errorf("cant parse jwt token: %v", err)
	}
	if payload.Valid() != nil {
		return false, fmt.Errorf("invalid jwt token: %v", err)
	}
	return payload.SessionID == s.ID && payload.UserID == s.UserID, nil
}

// Middleware для проверки CSRF
func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Try to get token from header first (for Swagger)
		if authHeader := r.Header.Get("X-Session-Token"); authHeader != "" {
			token = authHeader
		} else {
			// Try to get token from cookie (for browser)
			if cookie, err := r.Cookie("session_token"); err == nil {
				token = cookie.Value
			}
		}

		if r.Method != "OPTIONS" && r.Method != "GET" {
			sessionToken := token // Получить из сессии
			formToken := r.FormValue("csrf_token")

			if sessionToken != formToken || formToken == "" {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		// Генерация нового токена для формы
		newToken := generateCSRFToken()
		setSessionToken(w, r, newToken) // Сохранить в сессии

		next(w, r)
	}
}
