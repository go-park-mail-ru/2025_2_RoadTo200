package app

import (
	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	websocket "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/websocket"
)

func (a *App) initHandlers() {
	a.handlers = &Handlers{
		Auth:      handler.NewAuthHandler(a.services.Auth, a.logger),
		Session:   handler.NewSessionHandler(a.services.Auth, a.logger),
		Feed:      handler.NewFeedHandler(a.services.Feed, a.logger),
		Profile:   handler.NewProfileHandler(a.services.Profile, a.logger),
		Swipe:     handler.NewSwipeHandler(a.services.Swipe, a.logger),
		Match:     handler.NewMatchHandler(a.services.Match, a.logger),
		Chat:      handler.NewChatHandler(a.services.Chat, a.logger),
		WebSocket: websocket.NewWebSocketHandler(a.services.Chat, a.services.Auth, a.resources.RedisPubSub, a.logger),
		Strike:  handler.NewStrikeHandler(a.services.Strike, a.logger),
	}
}
