package app

import (
	"fmt"
	"net/http"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics/prometheus"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics/web"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *App) initServer() {
	a.server = httpserver.NewServer()

	// Settip metrics
	a.setupMetrics()

	// Global middleware
	a.server.AddMiddleware(middleware.LogMiddleware(a.logger))
	a.server.AddMiddleware(middleware.CORSMiddleware(&a.config.Cors))

	a.setupPublicRoutes()
	a.setupUtilRoutes()
	a.setupProtectedRoutes()
}

func (a *App) setupMetrics() {
	mv := web.NewHttpMetricCollector(prometheus.NewPrometheusHttpMetrics(a.config.Name, a.metrics)).Middleware()
	a.server.AddHandler("/metrics", promhttp.HandlerFor(a.metrics, promhttp.HandlerOpts{}))
	a.server.SetMetricMiddleware(mv)
}

func (a *App) setupUtilRoutes() {
	// Swagger should be accessible without auth
	a.server.AddHandler("/swagger/", http.HandlerFunc(handler.RegisterSwagger(a.config.SwaggerPath)))

	// Health check
	a.server.GET("/api/health", handler.HealthHandler)
}

func (a *App) setupPublicRoutes() {
	// Public routes (no auth required)
	a.server.POST("/api/register", a.handlers.Auth.Register)
	a.server.POST("/api/login", a.handlers.Auth.Login)
	a.server.GET("/api/session", a.handlers.Session.GetSession)

	// WebSocket routes (auth handled internally)
	a.server.GET("/ws/chat", a.handlers.WebSocket.HandleConnection)
	a.server.GET("/ws/notifications", a.handlers.NotificationWS.HandleConnection)

	// Payment webhook (no auth required, but signature is verified)
	a.server.POST("/api/notificate_premium", a.handlers.Payment.HandleWebhook)
}

// Protected routes (require auth)
func (a *App) setupProtectedRoutes() {
	// Auth middleware for protected routes
	a.server.SetAuthMiddleware(middleware.AuthMiddleware(a.services.Auth))

	// Auth endpoints
	a.server.POST("/api/logout", a.handlers.Auth.Logout)

	// Profile endpoints
	a.server.GET("/api/profile", a.handlers.Profile.GetProfile)
	a.server.GET("/api/profile/{id}", a.handlers.Profile.GetProfileByID)
	a.server.PUT("/api/profile/info", a.handlers.Profile.UpdateProfileInfo)
	a.server.PUT("/api/profile/preference", a.handlers.Profile.UpdatePreferences)
	a.server.PUT("/api/profile/interest", a.handlers.Profile.UpdateInterests)
	a.server.PUT("/api/profile/photo/{id}", a.handlers.Profile.SetPrimaryPhoto)
	a.server.DELETE("/api/profile/photo/{id}", a.handlers.Profile.DeletePhoto)
	a.server.POST("/api/profile/photo", a.handlers.Profile.UploadPhotos)

	// Feed endpoints
	a.server.GET("/api/feed", a.handlers.Feed.GetFeed)

	// Swipe endpoints
	a.server.POST("/api/swipe", a.handlers.Swipe.ProcessSwipe)

	// Match endpoints
	a.server.GET("/api/match", a.handlers.Match.GetUserMatches)
	a.server.DELETE("/api/match", a.handlers.Match.Unmatch)

	// Chat routes
	a.server.GET("/api/chats", a.handlers.Chat.GetConversations)
	a.server.GET("/api/chats/unread", a.handlers.Chat.GetUnreadCount)
	a.server.GET("/api/chats/{match_id}", a.handlers.Chat.GetMessages)
	a.server.POST("/api/chats/{match_id}/messages", a.handlers.Chat.SendMessage)
	a.server.POST("/api/chats/{match_id}/read", a.handlers.Chat.MarkAsRead)
	// Strike endpoints
	a.server.POST("/api/strike", a.handlers.Strike.CreateStrike)
	a.server.GET("/api/strike/{id}", a.handlers.Strike.GetStrike)
	a.server.PUT("/api/strike/{id}/status", a.handlers.Strike.UpdateStrikeStatus)
	a.server.DELETE("/api/strike/{id}", a.handlers.Strike.DeleteStrike)
	a.server.GET("/api/strike/user/{user_id}", a.handlers.Strike.GetStrikesByUserID)
	a.server.GET("/api/strike/type/{type}", a.handlers.Strike.GetStrikesByType)
	a.server.GET("/api/strike/range", a.handlers.Strike.GetStrikesByDateRange)
	a.server.GET("/api/strike/user/{user_id}/stat", a.handlers.Strike.GetUserStrikeStats)

	// Payment endpoints
	a.server.POST("/api/payment/create", a.handlers.Payment.CreatePayment)

	// Notification endpoints
	a.server.GET("/api/notifications", a.handlers.Notification.GetNotifications)
	a.server.PUT("/api/notifications/{notification_id}/read", a.handlers.Notification.MarkAsRead)
}

func (a *App) runServer() {
	addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	a.logger.Infof("Server starting on %s", addr)

	if err := a.server.Run(addr); err != nil {
		logger.Fatal("Server failed: ", err)
	}
	a.logger.Info("Server stopped")
}
