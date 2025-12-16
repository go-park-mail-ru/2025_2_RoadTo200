package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	profileService service.ProfileService
	logger         logger.Log
}

func NewProfileHandler(profileService service.ProfileService, l logger.Log) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         l,
	}
}

// GetProfile godoc
// @Summary Получить профиль пользователя
// @Description Возвращает полную информацию о профиле текущего пользователя
// @Tags profile
// @Produce json
// @Security SessionToken
// @Success 200 {object} dto.ProfileResponse "Профиль пользователя"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Профиль не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/profile [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.GetProfile")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, expectation.ErrProfileNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "profile not found")
			return
		}
		h.logger.Errorf("GetProfile: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := dto.ProfileResponse{
		User:        profile.User,
		Preferences: profile.Preferences,
		Photos:      profile.Photos,
		Interests:   profile.Interests,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetProfileByID godoc
// @Summary Получить профиль пользователя по ID
// @Description Возвращает полную информацию о профиле указанного пользователя
// @Tags profile
// @Produce json
// @Security SessionToken
// @Param id path string true "ID пользователя"
// @Success 200 {object} dto.ProfileResponse "Профиль пользователя"
// @Failure 400 {object} map[string]string "Неверный ID пользователя"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Профиль не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/profile/{id} [get]
func (h *ProfileHandler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.GetProfileByID")

	// Проверяем авторизацию (пользователь должен быть авторизован)
	_, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получаем ID пользователя из пути
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		// Если ID не указан в пути, пытаемся извлечь из URL
		path := r.URL.Path
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "profile" && i+1 < len(parts) {
				userIDStr = parts[i+1]
				break
			}
		}
	}

	if userIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warnf("Invalid user ID: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Получаем ID текущего пользователя (viewer)
	viewerID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получаем профиль пользователя с информацией об отношениях
	profile, err := h.profileService.GetProfileWithRelations(r.Context(), viewerID, userID)
	if err != nil {
		if errors.Is(err, expectation.ErrProfileNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "profile not found")
			return
		}
		h.logger.Errorf("GetProfileWithRelations: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := dto.ProfileResponse{
		User:        profile.User,
		Preferences: profile.Preferences,
		Photos:      profile.Photos,
		Interests:   profile.Interests,
		IsLiked:     profile.IsLiked,
		IsMatched:   profile.IsMatched,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// getContext функция извлечения контекста и проверки формата тела
func (h *ProfileHandler) getContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, errors.New("unauthorized")
	}

	// Только JSON запросы
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		h.logger.Warnf("getContext: Content-Type is not application/json")
		utils.WriteJSONError(w, http.StatusBadRequest, "only JSON content type supported for this endpoint")
		return uuid.Nil, errors.New("unauthorized")
	}

	return userID, nil
}

// UploadPhotos godoc
// @Summary Загрузить фотографии
// @Description Загрузка фотографий профиля
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security SessionToken
// @Param photos formData file true "Фотографии для загрузки"
// @Success 200 {object} domain.UserPhoto "Успешная загрузка фото"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Превышен лимит фото"
// @Router /api/profile/photo [post]
func (h *ProfileHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.UploadPhotos")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.logger.Infof("UploadPhotos request: userID=%s, Content-Type=%s, Content-Length=%s",
		userID, r.Header.Get("Content-Type"), r.Header.Get("Content-Length"))

	// Проверяем размер запроса
	if r.ContentLength > 0 {
		h.logger.Infof("Request content length: %d bytes (%.2f MB)",
			r.ContentLength, float64(r.ContentLength)/(1024*1024))
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Errorf("ParseMultipartForm failed: %v (max memory: %d bytes, content-length: %s)",
			err, 32<<20, r.Header.Get("Content-Length"))
		utils.WriteJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	h.logger.Infof("MultipartForm parsed successfully: form=%v", r.MultipartForm != nil)

	files := r.MultipartForm.File["photos"]
	h.logger.Infof("Files received: count=%d", len(files))

	if len(files) == 0 {
		h.logger.Warnf("uploadPhotos: no files uploaded")
		utils.WriteJSONError(w, http.StatusBadRequest, "no photos provided")
		return
	}

	// Логируем информацию о каждом файле
	for i, fileHeader := range files {
		h.logger.Infof("File %d: filename=%s, size=%d bytes (%.2f MB), Content-Type=%s",
			i+1, fileHeader.Filename, fileHeader.Size, float64(fileHeader.Size)/(1024*1024), fileHeader.Header.Get("Content-Type"))
	}

	uploadedPhotos, err := h.profileService.UploadPhotos(r.Context(), userID, files)
	if err != nil {
		h.logger.Errorf("uploadPhotos failed: %v", err)
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoLimitExceeded) {
			status = http.StatusForbidden
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	h.logger.Infof("UploadPhotos successful: uploaded %d photos", len(uploadedPhotos))
	utils.WriteJSON(w, http.StatusOK, uploadedPhotos)
}

// UpdateProfileInfo godoc
// @Summary Обновить основную информацию профиля
// @Description Обновляет основную информацию профиля пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body domain.ProfileUpdateRequest true "Данные для обновления профиля"
// @Success 200 {object} dto.SuccessResponse "Профиль успешно обновлен"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Профиль не найден"
// @Router /api/profile/info [put]
func (h *ProfileHandler) UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.UpdateProfileInfo")

	userID, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	var req domain.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdateProfileInfo(r.Context(), userID, &req); err != nil {
		h.logger.Errorf("updateProfileInfo: %v", err)
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrProfileNotFound) {
			status = http.StatusNotFound
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Profile updated successfully"})
}

// UpdatePreferences godoc
// @Summary Обновить предпочтения пользователя
// @Description Обновляет предпочтения пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body domain.PreferencesUpdateRequest true "Данные предпочтений"
// @Success 200 {object} dto.SuccessResponse "Предпочтения успешно обновлены"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Router /api/profile/preferences [put]
func (h *ProfileHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.UpdatePreferences")

	userID, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	var req domain.PreferencesUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdatePreferences(r.Context(), userID, &req); err != nil {
		h.logger.Errorf("updatePreferences: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Preferences updated successfully"})
}

// UpdateInterests godoc
// @Summary Обновить интересы пользователя
// @Description Обновляет интересы пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body []domain.Interest true "Список интересов"
// @Success 200 {object} dto.SuccessResponse "Интересы успешно обновлены"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Router /api/profile/interests [put]
func (h *ProfileHandler) UpdateInterests(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.UpdateInterests")

	userID, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	var req []domain.Interest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdateInterests(r.Context(), userID, req); err != nil {
		h.logger.Errorf("updatePreferences: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Preferences updated successfully"})
}

// HandlePhoto мультиплексер фото по методам
func (h *ProfileHandler) handlePhoto(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, error) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return uuid.Nil, uuid.Nil, errors.New("unauthorized")
	}

	photoID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Warnf("handlePhoto: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "photo ID is required")
		return uuid.Nil, uuid.Nil, errors.New("photo ID is required")
	}

	return userID, photoID, nil
}

// DeletePhoto godoc
// @Summary Удалить фотографию
// @Description Удаляет фотографию профиля
// @Tags profile
// @Produce json
// @Security SessionToken
// @Param id path string true "ID фотографии"
// @Success 200 {object} dto.SuccessResponse "Фото успешно удалено"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Фото не принадлежит пользователю"
// @Failure 404 {object} map[string]string "Фото не найдено"
// @Router /api/profile/photo/{id} [delete]
func (h *ProfileHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.DeletePhoto")

	userID, photoID, err := h.handlePhoto(w, r)
	if err != nil {
		return
	}

	if err := h.profileService.DeletePhoto(r.Context(), userID, photoID); err != nil {
		h.logger.Errorf("deletePhoto: %v", err)
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, expectation.ErrPhotoNotOwned) {
			status = http.StatusForbidden
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Photo deleted successfully"})
}

// SetPrimaryPhoto godoc
// @Summary Установить основное фото
// @Description Устанавливает основную фотографию профиля
// @Tags profile
// @Produce json
// @Security SessionToken
// @Param id path string true "ID фотографии"
// @Success 200 {object} dto.SuccessResponse "Основное фото успешно установлено"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Фото не принадлежит пользователю"
// @Failure 404 {object} map[string]string "Фото не найдено"
// @Router /api/profile/photo/{id} [put]
func (h *ProfileHandler) SetPrimaryPhoto(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("profileHandler.SetPrimaryPhoto")

	userID, photoID, err := h.handlePhoto(w, r)
	if err != nil {
		return
	}

	if err := h.profileService.SetPrimaryPhoto(r.Context(), userID, photoID); err != nil {
		h.logger.Errorf("setPrimaryPhoto: %v", err)
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, expectation.ErrPhotoNotOwned) {
			status = http.StatusForbidden
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Primary photo set successfully"})
}
