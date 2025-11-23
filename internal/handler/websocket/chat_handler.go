package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	chatService service.ChatService
	authService service.AuthService
	logger      logger.Log
}

func NewWebSocketHandler(chatService service.ChatService, authService service.AuthService, logger logger.Log) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
		authService: authService,
		logger:      logger,
	}
}

// HandleConnection handles WebSocket requests from the peer.
func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("WebSocketHandler.HandleConnection")

	// 1. Authenticate user
	token := r.URL.Query().Get("session_token")
	if token == "" {
		h.logger.Warn("No session token provided")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	h.logger.Infof("User %s connected via WebSocket", user.ID)

	// 3. Subscribe to Redis
	msgChan, cleanup, err := h.chatService.Subscribe(r.Context(), user.ID)
	if err != nil {
		h.logger.Errorf("Failed to subscribe to chat service: %v", err)
		return
	}
	defer cleanup()

	// 4. Start loops
	go h.writePump(conn, msgChan)
	h.readPump(conn, user.ID)
}

// readPump pumps messages from the websocket connection to the hub.
func (h *WebSocketHandler) readPump(conn *websocket.Conn, userID uuid.UUID) {
	defer func() {
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Errorf("error: %v", err)
			}
			break
		}

		// Parse message
		var req dto.SendMessageRequest
		if err := json.Unmarshal(message, &req); err != nil {
			h.logger.Warnf("Invalid message format: %v", err)
			continue
		}

		// Send message via service
		// Note: We use background context here because the request context might be closed if connection drops?
		// Actually, we should probably use a context that is tied to the connection lifecycle.
		// For now, using Background is okay for the service call, but we should handle timeouts.
		_, err = h.chatService.SendMessage(context.Background(), userID, &req)
		if err != nil {
			h.logger.Errorf("Failed to send message: %v", err)
			// Optionally send error back to client
			errMsg := domain.ChatMessage{
				Type:  "error",
				Error: "Failed to send message",
			}
			if data, err := json.Marshal(errMsg); err == nil {
				conn.WriteMessage(websocket.TextMessage, data)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (h *WebSocketHandler) writePump(conn *websocket.Conn, msgChan <-chan domain.ChatMessage) {
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
				// The hub closed the channel.
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				h.logger.Errorf("Failed to marshal message: %v", err)
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
