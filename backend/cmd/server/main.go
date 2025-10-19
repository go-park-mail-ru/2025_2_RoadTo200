package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
)

// @title Terabithia API
// @version 1.0
// @description API для dating приложения
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@datingapp.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey SessionToken
// @in cookie
// @name session_token
func main() {
	// Инициализация репозиториев
	userRepo := repository.NewInMemoryUserRepository()
	sessionRepo := repository.NewInMemorySessionRepository()

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, sessionRepo)

	// Инициализация обработчиков
	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(authService)

	// Middleware
	corsMiddleware := middleware.CORSMiddleware
	//authMiddleware := middleware.AuthMiddleware(authService)

	// Маршруты
	http.Handle("/api/register", corsMiddleware(http.HandlerFunc(authHandler.Register)))
	http.Handle("/api/login", corsMiddleware(http.HandlerFunc(authHandler.Login)))
	http.Handle("/api/session", corsMiddleware(http.HandlerFunc(sessionHandler.GetSession)))
	http.Handle("/api/logout", corsMiddleware(http.HandlerFunc(authHandler.Logout)))

	// Защищенные маршруты
	//http.Handle("/api/feed", corsMiddleware(authMiddleware(http.HandlerFunc(handler.FeedHandler))))
	//http.Handle("/api/swipe", corsMiddleware(authMiddleware(http.HandlerFunc(handler.SwipeHandler))))

	fmt.Println("Server running on http://:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
