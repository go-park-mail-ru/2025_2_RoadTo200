package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService service.ChatService
	logger      logger.Log
}

func NewChatHandler(chatService service.ChatService, logger logger.Log) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		logger:      logger,
	}
}

// GetConversations returns list of all chats with last message
// @Summary Get user conversations
// @Description Get list of all chats with last message and unread count
// @Tags chat
// @Accept json
// @Produce json
// @Param search query string false "Поиск чатов"
// @Success 200 {object} map[string]interface{} "conversations"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal error"
// @Router /api/chats [get]
func (h *ChatHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("ChatHandler.GetConversations")

	// Get user from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("GetUserIDFromContext failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get search query
	searchQuery := r.URL.Query().Get("search")

	// Get conversations
	conversations, err := h.chatService.GetConversations(r.Context(), userID, searchQuery)
	if err != nil {
		h.logger.Errorf("Failed to get conversations: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conversations": conversations,
	})
}

// GetMessages returns message history for a specific match
// @Summary Get chat messages
// @Description Get message history for a specific match with pagination
// @Tags chat
// @Accept json
// @Produce json
// @Param match_id path string true "Match ID"
// @Param limit query int false "Limit (default 50)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} map[string]interface{} "messages"
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "match not found"
// @Failure 500 {string} string "internal error"
// @Router /api/chats/{match_id} [get]
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("ChatHandler.GetMessages")

	// Get user from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("GetUserIDFromContext failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get match_id from URL
	matchIDStr := r.PathValue("match_id")
	if matchIDStr == "" {
		h.logger.Error("match_id is empty")
		http.Error(w, "match_id is required", http.StatusBadRequest)
		return
	}

	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		h.logger.Errorf("Invalid match_id: %v", err)
		http.Error(w, "invalid match_id", http.StatusBadRequest)
		return
	}

	// Get pagination params
	limit := 50 // default
	offset := 0 // default

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get messages
	messages, err := h.chatService.GetMessages(r.Context(), userID, matchID, limit, offset)
	if err != nil {
		if err == errors.ErrMatchNotFound {
			http.Error(w, "match not found", http.StatusNotFound)
			return
		}
		if err == errors.ErrNotMatchParticipant {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		h.logger.Errorf("Failed to get messages: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

// MarkAsRead marks all messages in a match as read
// @Summary Mark messages as read
// @Description Mark all messages in a specific match as read for the current user
// @Tags chat
// @Accept json
// @Produce json
// @Param match_id path string true "Match ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "match not found"
// @Failure 500 {string} string "internal error"
// @Router /api/chats/{match_id}/read [post]
func (h *ChatHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("ChatHandler.MarkAsRead")

	// Get user from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("GetUserIDFromContext failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get match_id from URL
	matchIDStr := r.PathValue("match_id")
	if matchIDStr == "" {
		h.logger.Error("match_id is empty")
		http.Error(w, "match_id is required", http.StatusBadRequest)
		return
	}

	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		h.logger.Errorf("Invalid match_id: %v", err)
		http.Error(w, "invalid match_id", http.StatusBadRequest)
		return
	}

	// Mark as read
	if err := h.chatService.MarkAsRead(r.Context(), userID, matchID); err != nil {
		if err == errors.ErrMatchNotFound {
			http.Error(w, "match not found", http.StatusNotFound)
			return
		}
		if err == errors.ErrNotMatchParticipant {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		h.logger.Errorf("Failed to mark as read: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "marked as read",
	})
}

// SendMessage sends a new message
// @Summary Send a message
// @Description Send a new message to a match
// @Tags chat
// @Accept json
// @Produce json
// @Param match_id path string true "Match ID"
// @Param input body dto.SendMessageRequest true "Message content"
// @Success 201 {object} map[string]interface{} "message_id, created_at"
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden or match not active"
// @Failure 404 {string} string "match not found"
// @Failure 500 {string} string "internal error"
// @Router /api/chats/{match_id}/messages [post]
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("ChatHandler.SendMessage")

	// Get user from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("GetUserIDFromContext failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get match_id from URL
	matchIDStr := r.PathValue("match_id")
	if matchIDStr == "" {
		h.logger.Error("match_id is empty")
		http.Error(w, "match_id is required", http.StatusBadRequest)
		return
	}

	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		h.logger.Errorf("Invalid match_id: %v", err)
		http.Error(w, "invalid match_id", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate content
	if len(req.Content) == 0 || len(req.Content) > 1000 {
		http.Error(w, "content must be between 1 and 1000 characters", http.StatusBadRequest)
		return
	}

	// Send message
	sendReq := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: req.Content,
	}

	message, err := h.chatService.SendMessage(r.Context(), userID, sendReq)
	if err != nil {
		if err == errors.ErrMatchNotFound {
			http.Error(w, "match not found", http.StatusNotFound)
			return
		}
		if err == errors.ErrMatchNotActive {
			http.Error(w, "match is not active", http.StatusForbidden)
			return
		}
		if err == errors.ErrNotMatchParticipant {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		h.logger.Errorf("Failed to send message: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message_id": message.ID,
		"created_at": message.CreatedAt,
	})
}

// GetUnreadCount returns total unread message count
// @Summary Get unread messages count
// @Description Get total count of unread messages for the current user
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int "unread_count"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal error"
// @Router /api/chats/unread [get]
func (h *ChatHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("ChatHandler.GetUnreadCount")

	// Get user from context
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("GetUserIDFromContext failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get unread count
	count, err := h.chatService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get unread count: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{
		"unread_count": count,
	})
}
