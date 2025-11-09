package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
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
// @Success 200 {object} ProfileResponse "Профиль пользователя"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Профиль не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/profile/profile [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		h.logger.Warnf("GetProfile: %v", err)
		if err == expectation.ErrProfileNotFound {
			utils.WriteJSONError(w, http.StatusNotFound, "profile not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := dto.ProfileResponse{
		User:        profile.User,
		Preferences: profile.Preferences,
		Photos:      profile.Photos,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// ChangeProfileJSON godoc
// @Summary Изменить профиль (JSON)
// @Description Изменение профиля через JSON для всех действий кроме загрузки фото
// @Description - updateInfo: обновление основной информации
// @Description - updatePreferences: обновление предпочтений
// @Tags profile
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body dto.UpdateProfileRequest true "Данные для изменения"
// @Success 200 {object} dto.SuccessResponse "Успешное обновление"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Router /api/profile/changeProfile [post]
func (h *ProfileHandler) getContext(r *http.Request) (uuid.UUID, error) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		//utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, errors.New("unauthorized")
	}

	// Только JSON запросы
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		//utils.WriteJSONError(w, http.StatusBadRequest, "only JSON content type supported for this endpoint")
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
// @Success 200 {object} dto.UploadPhotosResponse "Успешная загрузка фото"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Превышен лимит фото"
// @Router /api/profile/uploadPhotos [post]
func (h *ProfileHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Warnf("ParseMultipartForm: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		h.logger.Warnf("uploadPhotos: no files uploaded")
		utils.WriteJSONError(w, http.StatusBadRequest, "no photos provided")
		return
	}

	uploadedPhotos, err := h.profileService.UploadPhotos(userID, files)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoLimitExceeded) {
			status = http.StatusForbidden
		}
		h.logger.Warnf("uploadPhotos: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, uploadedPhotos)
}

// TODO: Swagger документация для метода
// updateProfileInfo обновляет основную информацию профиля
func (h *ProfileHandler) UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getContext(r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
	}

	var req domain.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdateProfileInfo(userID, &req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrProfileNotFound) {
			status = http.StatusNotFound
		}
		h.logger.Warnf("updateProfileInfo: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Profile updated successfully"})
}

// TODO: Swagger документация для метода
// updatePreferences обновляет предпочтения пользователя
func (h *ProfileHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getContext(r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
	}

	var req domain.PreferencesUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdatePreferences(userID, &req); err != nil {
		h.logger.Warnf("updatePreferences: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Preferences updated successfully"})
}

// TODO: Swagger документация для метода
// UpdateInterest обновляет предпочтения пользователя
func (h *ProfileHandler) UpdateInterests(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getContext(r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
	}

	var req []domain.Interest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	h.logger.Debugf("handleJSONRequest: %v", req)

	if err := h.profileService.UpdateInterests(userID, req); err != nil {
		h.logger.Warnf("updatePreferences: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Preferences updated successfully"})
}

func (h *ProfileHandler) HandlePhoto(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getContext(r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, err.Error())
	}
	photoID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Warnf("handlePhoto: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "photo ID is required")
	}

	if r.Method != http.MethodPut {
		h.setPrimaryPhoto(w, userID, photoID)
	} else if r.Method == http.MethodDelete {
		h.deletePhoto(w, userID, photoID)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// TODO: Swagger документация для метода
// deletePhoto удаляет фотографию
func (h *ProfileHandler) deletePhoto(w http.ResponseWriter, userID uuid.UUID, photoID uuid.UUID) {

	if err := h.profileService.DeletePhoto(userID, photoID); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, expectation.ErrPhotoNotOwned) {
			status = http.StatusForbidden
		}
		h.logger.Warnf("deletePhoto: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Photo deleted successfully"})
}

// TODO: Swagger документация для метода
// setPrimaryPhoto устанавливает основное фото
func (h *ProfileHandler) setPrimaryPhoto(w http.ResponseWriter, userID uuid.UUID, photoID uuid.UUID) {
	if err := h.profileService.SetPrimaryPhoto(userID, photoID); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrPhotoNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, expectation.ErrPhotoNotOwned) {
			status = http.StatusForbidden
		}
		h.logger.Warnf("setPrimaryPhoto: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Primary photo set successfully"})
}
