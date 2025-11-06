package main

import (
	"context"
	"fmt"
	"net/http"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/minio"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/redis"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	postgres_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redis_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"

	_ "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/api/docs" // импорт сгенерированной docs
)

// @title Terabithia API
// @version 1.0
// @description API для dating приложения
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@datingapp.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 217.16.17.116:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @schemes http
// @in header
// @name Authorization

// @securityDefinitions.apikey SessionToken
// @in cookie
// @name session_token
func main() {
	// Загрузка конфигурации
	logger.Println("🚀 Starting server...")
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("❌ Failed to load config: ", err)
	}
	logger.Printf("✅ Config loaded: %s:%d", cfg.Host, cfg.Port)

	logg := logger.New(&cfg.Logger)

	// Инициализация PostgreSQL через pkg/postgres
	logger.Println("🔌 Connecting to PostgreSQL...")
	pool, err := postgres_connect.NewConnect(context.Background(), &cfg.Postgres)
	if err != nil {
		logger.Fatal("❌ Failed to connect to PostgreSQL: ", err)
	}
	defer pool.Close()

	// Инициализация Redis через pkg/redis
	logger.Println("🔌 Connecting to Redis...")
	redisPool, err := redis_connect.NewConnection(&cfg.Redis)
	if err != nil {
		logger.Fatal("❌ Failed to connect to Redis: ", err)
	}
	defer redisPool.Close()

	// Инициализация MinIO
	minioStorage, err := minio.NewMinIOStorage(&cfg.MinIO)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO: ", err)
	}
	logger.Println("✅ MinIO connected successfully")

	logger.Println("✅ All connections established")

	// Репозитории
	userRepo := postgres.NewUserRepository(pool)
	sessionRepo := redis.NewSessionRepository(redisPool)

	// Сервисы
	authService := service.NewAuthService(userRepo, sessionRepo, logg)
	feedService := service.NewFeedService(userRepo)
	profileService := service.NewProfileService(
		userRepo,
		postgres.NewUserPhotoRepository(pool),      // нужно создать эту реализацию
		postgres.NewUserPreferenceRepository(pool), // и эту
		minioStorage,
	)
	swipeService := service.NewSwipeService(
		postgres.NewSwipeRepository(pool),
		postgres.NewMatchRepository(pool),
	)
	matchService := service.NewMatchService(
		postgres.NewMatchRepository(pool),
		userRepo,
		postgres.NewSwipeRepository(pool),
	)

	// Обработчики
	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(authService)
	feedHandler := handler.NewFeedHandler(feedService)
	profileHandler := handler.NewProfileHandler(profileService)
	swipeHandler := handler.NewSwipeHandler(swipeService)
	matchHandler := handler.NewMatchHandler(matchService)

	server := httpserver.NewServer()
	// Middleware
	server.AddMiddleware(middleware.LogMiddleware(logg))
	server.AddMiddleware(middleware.CORSMiddleware(&cfg.Cors))

	server.AddHandler("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://217.16.17.116:%d/swagger/doc.json", cfg.Port)), // URL для doc.json
	))

	// Маршруты
	server.AddHandler("/api/register", http.HandlerFunc(authHandler.Register))
	server.AddHandler("/api/login", http.HandlerFunc(authHandler.Login))
	server.AddHandler("/api/session", http.HandlerFunc(sessionHandler.GetSession))
	server.AddHandler("/api/logout", http.HandlerFunc(authHandler.Logout))

	server.AddMiddleware(middleware.AuthMiddleware(authService))
	// Защищенные маршруты
	server.AddHandler("/api/profile/profile", http.HandlerFunc(profileHandler.GetProfile))
	server.AddHandler("/api/profile/changeProfile", http.HandlerFunc(profileHandler.ChangeProfile))
	server.AddHandler("/api/feed", http.HandlerFunc(feedHandler.GetFeed))
	server.AddHandler("/api/swipe", http.HandlerFunc(swipeHandler.ProcessSwipe))
	server.AddHandler("/api/matches", http.HandlerFunc(matchHandler.GetUserMatches))
	server.AddHandler("/api/matches/unmatch", http.HandlerFunc(matchHandler.Unmatch))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Printf("Server starting on %s", addr)

	if err := server.Run(addr); err != nil {
		logger.Fatal("Server failed: ", err)
	}
}

// Health handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"app":    "terabithia app",
	})
}
