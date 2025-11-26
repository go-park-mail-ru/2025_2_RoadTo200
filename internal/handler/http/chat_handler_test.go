package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChatHandler_SendMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockChatService := mocks.NewMockChatService(ctrl)
	handler := NewChatHandler(mockChatService, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()

	messageReq := dto.SendMessageRequest{
		Content: "Hello!",
	}

	message := &domain.Message{
		ID:        uuid.New(),
		MatchID:   matchID,
		SenderID:  userID,
		Content:   "Hello!",
		CreatedAt: time.Now(),
	}

	mockChatService.EXPECT().
		SendMessage(gomock.Any(), userID, gomock.Any()).
		Return(message, nil)

	body, _ := json.Marshal(messageReq)
	req := httptest.NewRequest(http.MethodPost, "/api/chat/"+matchID.String()+"/messages", bytes.NewReader(body))
	req.SetPathValue("match_id", matchID.String())
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.SendMessage(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestChatHandler_SendMessage_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockChatService := mocks.NewMockChatService(ctrl)
	handler := NewChatHandler(mockChatService, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/api/chat/"+matchID.String()+"/messages", bytes.NewReader([]byte("invalid")))
	req.SetPathValue("match_id", matchID.String())
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.SendMessage(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatHandler_SendMessage_MatchNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockChatService := mocks.NewMockChatService(ctrl)
	handler := NewChatHandler(mockChatService, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()

	messageReq := dto.SendMessageRequest{
		Content: "Hello!",
	}

	mockChatService.EXPECT().
		SendMessage(gomock.Any(), userID, gomock.Any()).
		Return(nil, errors.ErrMatchNotFound)

	body, _ := json.Marshal(messageReq)
	req := httptest.NewRequest(http.MethodPost, "/api/chat/"+matchID.String()+"/messages", bytes.NewReader(body))
	req.SetPathValue("match_id", matchID.String())
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.SendMessage(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestChatHandler_GetMessages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockChatService := mocks.NewMockChatService(ctrl)
	handler := NewChatHandler(mockChatService, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()

	messages := []domain.Message{
		{ID: uuid.New(), Content: "Hello"},
	}

	mockChatService.EXPECT().
		GetMessages(gomock.Any(), userID, matchID, 50, 0).
		Return(messages, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/chat/"+matchID.String()+"/messages", nil)
	req.SetPathValue("match_id", matchID.String())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetMessages(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChatHandler_GetConversations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockChatService := mocks.NewMockChatService(ctrl)
	handler := NewChatHandler(mockChatService, mockLogger)

	userID := uuid.New()

	conversations := []domain.Conversation{
		{MatchID: uuid.New()},
	}

	mockChatService.EXPECT().
		GetConversations(gomock.Any(), userID, "").
		Return(conversations, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/chat/conversations", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetConversations(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
