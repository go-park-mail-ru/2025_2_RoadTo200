package app

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/websocket"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
	"github.com/gomodule/redigo/redis"
	"github.com/minio/minio-go/v7"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	config       *config.Config
	logger       logger.Log
	server       *httpserver.Server
	metrics      *Metrics
	resources    *Resources
	repositories *Repositories
	services     *Services
	handlers     *Handlers
}

type Metrics struct {
	HttpMetrics metrics.HttpMetrics
	GrpcMetrics metrics.GrpcMetrics
}

type Resources struct {
	Postgres    interfaces.PgxIface
	Redis       *redis.Pool
	RedisPubSub *goredis.Client // for Pub/Sub (different from session Redis pool)
	MinIO       *minio.Client
}

type Repositories struct {
	Storage      interfaces.FileStorage
	Match        interfaces.MatchRepository
	Message      interfaces.MessageRepository
	Subscription interfaces.SubscriptionRepository
	Swipe        interfaces.SwipeRepository
	Preference   interfaces.UserPreferenceRepository
	Photo        interfaces.UserPhotoRepository
	Session      interfaces.SessionRepository
	User         interfaces.UserRepository
	Strike       interfaces.StrikeRepository
}

type Services struct {
	Auth    service.AuthService
	Feed    service.FeedService
	Profile service.ProfileService
	Swipe   service.SwipeService
	Match   service.MatchService
	Chat    service.ChatService
	Strike  service.StrikeService
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
	if err := app.initResources(); err != nil {
		app.logger.Fatal(err)
		return
	}
	if err := app.runMigrations(); err != nil {
		app.logger.Fatal(err)
		return
	}

	app.RegisterMetrics()
	app.initRepository()
	app.logger.Info("✅ Repositories initialized")
	err := app.initServices()
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
