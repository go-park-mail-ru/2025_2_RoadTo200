package app

import (
	"fmt"
	"net/http"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
)

func (a *App) initServer() {
	a.server = httpserver.NewServer()

	// Global middleware
	a.server.AddMiddleware(middleware.LogMiddleware(a.logger))
	a.server.AddMiddleware(middleware.CORSMiddleware(&a.config.Cors))

	a.setupPublicRoutes()
	a.setupUtilRoutes()
	a.setupProtectedRoutes()
}

func (a *App) setupUtilRoutes() {
	// Swagger should be accessible without auth
	a.server.AddHandler("/swagger/", http.HandlerFunc(handler.SwaggerHandler))

	// Health check
	a.server.GET("/api/health", handler.HealthHandler)
}

func (a *App) setupPublicRoutes() {
	// Public routes (no auth required)
	a.server.POST("/api/register", a.handlers.Auth.Register)
	a.server.POST("/api/login", a.handlers.Auth.Login)
	a.server.GET("/api/session", a.handlers.Session.GetSession)
}

func (a *App) setupProtectedRoutes() {
	// Auth middleware for protected routes
	a.server.AddMiddleware(middleware.AuthMiddleware(a.services.Auth))

	// Protected routes (require auth)
	a.server.POST("/api/logout", a.handlers.Auth.Logout)
	a.server.GET("/api/profile", a.handlers.Profile.GetProfile)
	a.server.PUT("/api/profile/info", a.handlers.Profile.UpdateProfileInfo)
	a.server.PUT("/api/profile/preference", a.handlers.Profile.UpdatePreferences)
	a.server.PUT("/api/profile/interest", a.handlers.Profile.UpdateInterests)
	a.server.PUT("/api/profile/photo/{id}", a.handlers.Profile.SetPrimaryPhoto)
	a.server.DELETE("/api/profile/photo/{id}", a.handlers.Profile.DeletePhoto)
	a.server.POST("/api/profile/photo", a.handlers.Profile.UploadPhotos)
	a.server.GET("/api/feed", a.handlers.Feed.GetFeed)
	a.server.POST("/api/swipe", a.handlers.Swipe.ProcessSwipe)
	a.server.GET("/api/match", a.handlers.Match.GetUserMatches)
	a.server.DELETE("/api/match", a.handlers.Match.Unmatch)
	a.server.POST("/api/report", a.handlers.Report.CreateSupportTicket)
	a.server.GET("/api/report", a.handlers.Report.GetUserSupportTickets)
	a.server.GET("/api/report/{id}", a.handlers.Report.GetUserSupportTicket)
	a.server.GET("/api/report/stats", a.handlers.Report.GetSupportStats)
}

func (a *App) runServer() {
	addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	a.logger.Infof("Server starting on %s", addr)

	if err := a.server.Run(addr); err != nil {
		logger.Fatal("Server failed: ", err)
	}
	a.logger.Info("Server stopped")
}
