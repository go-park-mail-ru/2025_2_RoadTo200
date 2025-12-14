package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type NotificationWebSocketHandler struct {
	authService service.AuthService
	redisClient *redis.Client
	logger      logger.Log
}

func NewNotificationWebSocketHandler(authService service.AuthService, redisClient *redis.Client, logger logger.Log) *NotificationWebSocketHandler {
	return &NotificationWebSocketHandler{
		authService: authService,
		redisClient: redisClient,
		logger:      logger,
	}
}

// HandleConnection handles WebSocket requests for notifications
func (h *NotificationWebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("NotificationWebSocketHandler.HandleConnection")

	// 1. Authenticate user
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		h.logger.Warn("No session token provided")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := cookie.Value

	user, err := h.authService.ValidateSession(r.Context(), token)
	if err != nil {
		h.logger.Warnf("Invalid session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	h.logger.Infof("User %s connected via WebSocket for notifications", user.ID)

	// 3. Subscribe to Redis notifications channel
	msgChan, cleanup, err := h.subscribe(r.Context(), user.ID)
	if err != nil {
		h.logger.Errorf("Failed to subscribe to notifications: %v", err)
		return
	}
	defer cleanup()

	// 4. Start loops
	go h.writePump(conn, msgChan)
	h.readPump(conn)
}

// subscribe returns a channel that receives notifications for the user
func (h *NotificationWebSocketHandler) subscribe(ctx context.Context, userID uuid.UUID) (<-chan domain.NotificationMessage, func(), error) {
	notifChannel := fmt.Sprintf("notification:user:%s", userID.String())

	pubsub := h.redisClient.Subscribe(ctx, notifChannel)

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to redis: %w", err)
	}

	// Create a channel for messages
	msgChan := make(chan domain.NotificationMessage)

	// Start a goroutine to pump messages
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var notifMsg domain.NotificationMessage
			if err := json.Unmarshal([]byte(msg.Payload), &notifMsg); err != nil {
				h.logger.Errorf("Failed to unmarshal notification message: %v", err)
				continue
			}
			msgChan <- notifMsg
		}
		close(msgChan)
	}()

	// Return cleanup function
	cleanup := func() {
		pubsub.Close()
	}

	return msgChan, cleanup, nil
}

// readPump pumps messages from the websocket connection (for pings/pongs only)
func (h *NotificationWebSocketHandler) readPump(conn *websocket.Conn) {
	defer conn.Close()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Errorf("error: %v", err)
			}
			break
		}
		// Notifications are read-only, we don't process incoming messages
	}
}

// writePump pumps notifications to the websocket connection
func (h *NotificationWebSocketHandler) writePump(conn *websocket.Conn, msgChan <-chan domain.NotificationMessage) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-msgChan:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The channel closed
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				h.logger.Errorf("Failed to marshal notification: %v", err)
				return
			}
			w.Write(data)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
