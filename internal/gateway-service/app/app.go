package app

import (
	auth "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/service/interfaces"
	chat "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	core "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/service/interfaces"
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/websocket"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
	"github.com/prometheus/client_golang/prometheus"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	config      *config.Config
	logger      logger.Log
	server      *httpserver.Server
	metrics     *prometheus.Registry
	RedisPubSub *goredis.Client // for Pub/Sub (different from session Redis pool)
	services    *Services
	handlers    *Handlers
}

type Services struct {
	Auth    auth.AuthService
	Feed    core.FeedService
	Profile core.ProfileService
	Swipe   core.SwipeService
	Match   chat.MatchService
	Chat    chat.ChatService
	Strike  core.StrikeService
}

type Handlers struct {
	Auth      *handler.AuthHandler
	Session   *handler.SessionHandler
	Feed      *handler.FeedHandler
	Profile   *handler.ProfileHandler
	Swipe     *handler.SwipeHandler
	Match     *handler.MatchHandler
	Chat      *handler.ChatHandler
	WebSocket *websocket.WebSocketHandler
	Strike    *handler.StrikeHandler
}

func Run() {
	app := &App{}

	if err := app.initConfig(); err != nil {
		app.logger.Fatal(err)
		return
	}
	if err := app.initLogger(); err != nil {
		app.logger.Fatal(err)
		return
	}
	pub, err := redis.NewPubSubClient(&app.config.Redis)
	if err != nil {
		app.logger.Fatalf("Redis error %s", err)
		return
	}
	app.RedisPubSub = pub

	app.registerMetrics()
	app.logger.Info("✅ Repositories initialized")
	err = app.initServices()
	if err != nil {
		app.logger.Fatal(err)
	}
	app.logger.Info("✅ Services initialized")
	app.initHandlers()
	app.logger.Info("✅ Handlers initialized")
	app.initServer()
	app.logger.Info("✅ Server initialized")
	app.runServer()
}
