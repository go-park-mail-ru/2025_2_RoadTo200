package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/minio"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/redis"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	minio_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/minio"
	postgres_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redis_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"

	_ "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/api/docs"
)

//тест

// @title Terabithia Dating App API
// @version 1.0
// @description API для dating приложения Terabithia
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@terabithia.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 217.16.17.116:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey SessionToken
// @in header
// @name X-Session-Token
// @description Токен сессии для аутентификации пользователя (альтернатива cookie)

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
	minioPool, err := minio_connect.NewMinioPool(&cfg.MinIO)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO: ", err)
	}
	logger.Println("✅ MinIO connected successfully")

	logger.Println("✅ All connections established")

	// Репозитории
	storageRepo := minio.NewStorageRepository(minioPool, &cfg.MinIO)
	userRepo := postgres.NewUserRepository(pool, logg)
	sessionRepo := redis.NewSessionRepository(redisPool)
	preferenceRepo := postgres.NewUserPreferenceRepository(pool)
	photoRepo := postgres.NewUserPhotoRepository(pool)
	swipeRepo := postgres.NewSwipeRepository(pool)
	matchRepo := postgres.NewMatchRepository(pool)

	// Сервисы
	authService := service.NewAuthService(userRepo, sessionRepo, logg)
	feedService := service.NewFeedService(userRepo, preferenceRepo, photoRepo, logg)
	profileService := service.NewProfileService(
		userRepo,
		photoRepo,
		preferenceRepo,
		storageRepo,
	)
	swipeService := service.NewSwipeService(
		swipeRepo,
		matchRepo,
	)
	matchService := service.NewMatchService(
		matchRepo,
		userRepo,
		swipeRepo,
		photoRepo, // Добавляем photoRepo
	)

	// Обработчики
	authHandler := handler.NewAuthHandler(authService, logg)
	sessionHandler := handler.NewSessionHandler(authService, logg)
	feedHandler := handler.NewFeedHandler(feedService, logg)
	profileHandler := handler.NewProfileHandler(profileService, logg)
	swipeHandler := handler.NewSwipeHandler(swipeService)
	matchHandler := handler.NewMatchHandler(matchService)

	// TODO: Отдельный файл для хендлеров
	server := httpserver.NewServer()

	// Global middleware
	server.AddMiddleware(middleware.LogMiddleware(logg))
	server.AddMiddleware(middleware.CORSMiddleware(&cfg.Cors))

	// Public routes (no auth required)
	server.AddHandler("/api/register", http.HandlerFunc(authHandler.Register))
	server.AddHandler("/api/login", http.HandlerFunc(authHandler.Login))
	server.AddHandler("/api/session", http.HandlerFunc(sessionHandler.GetSession))

	// Swagger should be accessible without auth
	server.AddHandler("/swagger/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/doc.json" {
			data, err := os.ReadFile(getSwaggerPath())
			if err != nil {
				log.Printf("Error reading swagger.json: %v", err)
				utils.WriteJSONError(w, http.StatusInternalServerError, "Swagger docs not found: "+err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")).ServeHTTP(w, r)
	}))

	// Health check
	server.AddHandler("/health", http.HandlerFunc(healthHandler))

	// Auth middleware for protected routes
	server.AddMiddleware(middleware.AuthMiddleware(authService))

	// Protected routes (require auth)
	server.AddHandler("/api/logout", http.HandlerFunc(authHandler.Logout))
	server.AddHandler("/api/profile/profile", http.HandlerFunc(profileHandler.GetProfile))
	server.AddHandler("/api/profile/changeProfile", http.HandlerFunc(profileHandler.UpdateProfileInfo))    // JSON only
	server.AddHandler("/api/profile/changePreference", http.HandlerFunc(profileHandler.UpdatePreferences)) // JSON only
	server.AddHandler("/api/profile/changeInterest", http.HandlerFunc(profileHandler.UpdateInterests))     // JSON only
	server.AddHandler("/api/profile/photo/{id}", http.HandlerFunc(profileHandler.HandlePhoto))             // JSON only
	server.AddHandler("/api/profile/uploadPhotos", http.HandlerFunc(profileHandler.UploadPhotos))          // Multipart only
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

// Добавьте функцию для получения абсолютного пути
func getSwaggerPath() string {
	// Пробуем несколько возможных путей
	paths := []string{
		"docs/swagger.json",
		"./docs/swagger.json",
		"../docs/swagger.json",
		"../../docs/swagger.json",
	}

	// Получаем директорию где запущен бинарник
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths, filepath.Join(exeDir, "docs/swagger.json"))
	}

	// Текущая рабочая директория
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(wd, "docs/swagger.json"))
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Found swagger.json at: %s", path)
			return path
		}
	}

	return "docs/swagger.json" // fallback
}
