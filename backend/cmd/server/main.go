package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/redis"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	postgres_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redis_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
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
	// Загрузка конфигурации
	log.Println("🚀 Starting server...")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config: ", err)
	}
	log.Printf("✅ Config loaded: %s:%d", cfg.Host, cfg.Port)

	// Инициализация PostgreSQL через pkg/postgres
	log.Println("🔌 Connecting to PostgreSQL...")
	pool, err := postgres_connect.NewConnect(context.Background(), &cfg.Postgres)
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL: ", err)
	}
	defer pool.Close()

	// Инициализация Redis через pkg/redis
	log.Println("🔌 Connecting to Redis...")
	redisPool, err := redis_connect.NewConnection(&cfg.Redis)
	if err != nil {
		log.Fatal("❌ Failed to connect to Redis: ", err)
	}
	defer redisPool.Close()

	log.Println("✅ All connections established")
	// Репозитории
	userRepo := postgres.NewUserRepository(pool)
	sessionRepo := redis.NewSessionRepository(redisPool)

	// Сервисы
	authService := service.NewAuthService(userRepo, sessionRepo)

	// Обработчики
	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(authService)

	// Middleware
	corsMiddleware := middleware.CORSMiddleware

	// Маршруты
	http.Handle("/api/register", corsMiddleware(http.HandlerFunc(authHandler.Register)))
	http.Handle("/api/login", corsMiddleware(http.HandlerFunc(authHandler.Login)))
	http.Handle("/api/session", corsMiddleware(http.HandlerFunc(sessionHandler.GetSession)))
	http.Handle("/api/logout", corsMiddleware(http.HandlerFunc(authHandler.Logout)))

	// Защищенные маршруты (добавятся позже)
	// http.Handle("/api/feed", corsMiddleware(authMiddleware(http.HandlerFunc(feedHandler.Feed))))
	// http.Handle("/api/swipe", corsMiddleware(authMiddleware(http.HandlerFunc(swipeHandler.Swipe))))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
