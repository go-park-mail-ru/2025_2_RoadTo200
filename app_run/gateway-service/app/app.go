package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/websocket"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
	"github.com/prometheus/client_golang/prometheus"
	goredis "github.com/redis/go-redis/v9"
)

type Repositories struct {
	// Add repository fields if needed
}

type App struct {
	config       *config.Config
	logger       logger.Log
	server       *httpserver.Server
	metrics      *prometheus.Registry
	resources    *goredis.Client
	repositories *Repositories
	services     *Services
	handlers     *Handlers
}

type Services struct {
	Auth         service.AuthService
	Feed         service.FeedService
	Profile      service.ProfileService
	Swipe        service.SwipeService
	Match        service.MatchService
	Chat         service.ChatService
	Strike       service.StrikeService
	Payment      service.PaymentService
	Notification service.NotificationService
}

type Handlers struct {
	Auth           *handler.AuthHandler
	Session        *handler.SessionHandler
	Feed           *handler.FeedHandler
	Profile        *handler.ProfileHandler
	Swipe          *handler.SwipeHandler
	Match          *handler.MatchHandler
	Chat           *handler.ChatHandler
	WebSocket      *websocket.WebSocketHandler
	NotificationWS *websocket.NotificationWebSocketHandler
	Strike         *handler.StrikeHandler
	Payment        *handler.PaymentHandler
	Notification   *handler.NotificationHandler
}

func Run() {
	app := &App{}

	if err := app.initConfig(); err != nil {
		fmt.Printf("Failed to init config: %v\n", err)
		return
	}
	if err := app.initLogger(); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		return
	}
	if err := app.runMigrations(); err != nil {
		app.logger.Fatal(err)
		return
	}

	app.registerMetrics()
	app.logger.Info("✅ Metrics registered")

	rds, err := redis.NewPubSubClient(&app.config.Redis)
	if err != nil {
		app.logger.Fatal(err)
		return
	}
	app.resources = rds
	app.logger.Info("✅ Resources initialized")

	app.repositories = &Repositories{}
	app.logger.Info("✅ Repositories initialized")

	err = app.initServices()
	if err != nil {
		app.logger.Fatal(err)
		return
	}
	app.logger.Info("✅ Services initialized")
	app.initHandlers()
	app.logger.Info("✅ Handlers initialized")
	app.initServer()
	app.logger.Info("✅ Server initialized")
	app.runServer()
}
