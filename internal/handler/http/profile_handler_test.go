package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProfileHandler_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	profile := &domain.ProfileResponse{
		User: &domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		},
		Preferences: &domain.UserPreference{
			UserID: userID,
			AgeMin: 18,
			AgeMax: 35,
		},
		Photos:    []domain.UserPhoto{},
		Interests: []domain.Interest{},
	}

	mockProfileService.EXPECT().
		GetProfile(gomock.Any(), userID).
		Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetProfile(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.ProfileResponse
	// Используем стандартный json для десериализации
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Проверяем структуру ответа
	assert.NotNil(t, response.User)
	if user, ok := response.User.(*domain.User); ok {
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
	} else if userMap, ok := response.User.(map[string]interface{}); ok {
		// Если это map (из стандартного json), проверяем так
		assert.Equal(t, userID.String(), userMap["id"])
		assert.Equal(t, "test@example.com", userMap["email"])
	}
}

func TestProfileHandler_GetProfile_NoUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	// No userID in context
	rec := httptest.NewRecorder()

	handler.GetProfile(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProfileHandler_GetProfile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	mockProfileService.EXPECT().
		GetProfile(gomock.Any(), userID).
		Return(nil, errors.ErrProfileNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetProfile(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProfileHandler_UpdateProfileInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	updateData := map[string]interface{}{
		"name": "Updated Name",
		"bio":  "Updated Bio",
	}
	body, _ := json.Marshal(updateData)

	mockProfileService.EXPECT().
		UpdateProfileInfo(gomock.Any(), userID, gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile/info", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UpdateProfileInfo(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileHandler_UpdateProfileInfo_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/api/profile/info", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UpdateProfileInfo(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProfileHandler_UpdatePreferences_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	preferencesData := map[string]interface{}{
		"age_min":     18,
		"age_max":     35,
		"show_gender": "both",
	}
	body, _ := json.Marshal(preferencesData)

	mockProfileService.EXPECT().
		UpdatePreferences(gomock.Any(), userID, gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile/preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileHandler_UpdateInterests_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	interests := []domain.Interest{
		{
			UserID: userID,
			Theme:  "workout",
		},
	}
	body, _ := json.Marshal(interests)

	mockProfileService.EXPECT().
		UpdateInterests(gomock.Any(), userID, gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile/interests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UpdateInterests(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileHandler_UploadPhotos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("photos", "test.jpg")
	part.Write([]byte("test content"))
	writer.Close()

	uploadedPhotos := []domain.UserPhoto{
		{
			ID:       uuid.New(),
			UserID:   userID,
			PhotoURL: "http://example.com/test.jpg",
		},
	}

	mockProfileService.EXPECT().
		UploadPhotos(gomock.Any(), userID, gomock.Any()).
		Return(uploadedPhotos, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/profile/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UploadPhotos(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileHandler_UploadPhotos_NoFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/profile/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.UploadPhotos(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProfileHandler_DeletePhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()
	photoID := uuid.New()

	mockProfileService.EXPECT().
		DeletePhoto(gomock.Any(), userID, photoID).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/profile/photo/"+photoID.String(), nil)
	req.SetPathValue("id", photoID.String())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.DeletePhoto(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileHandler_DeletePhoto_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()
	photoID := uuid.New()

	mockProfileService.EXPECT().
		DeletePhoto(gomock.Any(), userID, photoID).
		Return(errors.ErrPhotoNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/profile/photo/"+photoID.String(), nil)
	req.SetPathValue("id", photoID.String())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.DeletePhoto(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProfileHandler_SetPrimaryPhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockProfileService := mocks.NewMockProfileService(ctrl)
	handler := NewProfileHandler(mockProfileService, mockLogger)

	userID := uuid.New()
	photoID := uuid.New()

	mockProfileService.EXPECT().
		SetPrimaryPhoto(gomock.Any(), userID, photoID).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile/photo/"+photoID.String(), nil)
	req.SetPathValue("id", photoID.String())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.SetPrimaryPhoto(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
