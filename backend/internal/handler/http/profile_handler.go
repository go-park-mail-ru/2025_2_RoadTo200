package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
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

// ProfileResponse represents profile response
type ProfileResponse struct {
	User        interface{} `json:"user"`
	Preferences interface{} `json:"preferences,omitempty"`
	Photos      interface{} `json:"photos,omitempty"`
}

// UpdateProfileInfoRequest запрос на обновление информации
type UpdateProfileInfoRequest struct {
	Action    string           `json:"action" example:"updateInfo"`
	Name      string           `json:"name,omitempty" example:"Алексей"`
	Phone     *string          `json:"phone,omitempty" example:"+79991234567"`
	BirthDate *time.Time       `json:"birth_date,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender    constants.Gender `json:"gender,omitempty" example:"male"`
	Bio       *string          `json:"bio,omitempty" example:"Люблю путешествия и спорт"`
	Latitude  *float64         `json:"latitude,omitempty" example:"55.7558"`
	Longitude *float64         `json:"longitude,omitempty" example:"37.6173"`
}

// UpdatePreferencesRequest запрос на обновление предпочтений
type UpdatePreferencesRequest struct {
	Action       string                     `json:"action" example:"updatePreferences"`
	ShowGender   constants.GenderPreference `json:"show_gender,omitempty" example:"female"`
	AgeMin       int                        `json:"age_min,omitempty" example:"18"`
	AgeMax       int                        `json:"age_max,omitempty" example:"35"`
	MaxDistance  int                        `json:"max_distance,omitempty" example:"50"`
	GlobalSearch bool                       `json:"global_search,omitempty" example:"false"`
}

// DeletePhotoRequest запрос на удаление фото
type DeletePhotoRequest struct {
	Action  string    `json:"action" example:"deletePhoto"`
	PhotoID uuid.UUID `json:"photo_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// SetPrimaryPhotoRequest запрос на установку основного фото
type SetPrimaryPhotoRequest struct {
	Action  string    `json:"action" example:"setPrimaryPhoto"`
	PhotoID uuid.UUID `json:"photo_id" example:"550e8400-e29b-41d4-a716-446655440000"`
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
	userID, err := h.getUserIDFromContext(r)
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

	response := service.ProfileResponse{
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
// @Param request body service.UpdateProfileRequest true "Данные для изменения"
// @Success 200 {object} service.SuccessResponse "Успешное обновление"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Router /api/profile/changeProfile [post]
func (h *ProfileHandler) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Только JSON запросы
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		utils.WriteJSONError(w, http.StatusBadRequest, "only JSON content type supported for this endpoint")
		return
	}

	h.handleJSONRequest(w, r, userID)
}

// UploadPhotos godoc
// @Summary Загрузить фотографии
// @Description Загрузка фотографий профиля
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security SessionToken
// @Param photos formData file true "Фотографии для загрузки"
// @Success 200 {object} service.UploadPhotosResponse "Успешная загрузка фото"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Превышен лимит фото"
// @Router /api/profile/uploadPhotos [post]
func (h *ProfileHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
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

	h.uploadPhotos(w, r, userID)
}

// handleJSONRequest обрабатывает JSON запросы
func (h *ProfileHandler) handleJSONRequest(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var req service.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("handleJSONRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	h.logger.Debugf("handleJSONRequest: %v", req)
	switch req.Action {
	case "updateInfo":
		h.updateProfileInfo(w, userID, req)
	case "updatePreferences":
		h.updatePreferences(w, userID, req)
	case "deletePhoto":
		h.deletePhoto(w, userID, req)
	case "setPrimaryPhoto":
		h.setPrimaryPhoto(w, userID, req)
	default:
		utils.WriteJSONError(w, http.StatusBadRequest, "unknown action")
	}
}

// handleMultipartRequest обрабатывает multipart запросы (загрузка фото)
func (h *ProfileHandler) handleMultipartRequest(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	fmt.Printf("Handling multipart request\n")
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		h.logger.Warnf("handleMultipartRequest: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	action := r.FormValue("action")
	if action != "uploadPhotos" {
		h.logger.Warnf("handleMultipartRequest: invalid action: %v", action)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid action for multipart request")
		return
	}

	h.uploadPhotos(w, r, userID)
}

// updateProfileInfo обновляет основную информацию профиля
func (h *ProfileHandler) updateProfileInfo(w http.ResponseWriter, userID uuid.UUID, req service.UpdateProfileRequest) {
	updateData := domain.ProfileUpdateRequest{
		Name:      req.Name,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Bio:       req.Bio,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// Обработка даты рождения
	if req.BirthDate != nil {
		updateData.BirthDate = req.BirthDate
	}

	if err := h.profileService.UpdateProfileInfo(userID, &updateData); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, expectation.ErrProfileNotFound) {
			status = http.StatusNotFound
		}
		h.logger.Warnf("updateProfileInfo: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, service.SuccessResponse{Message: "Profile updated successfully"})
}

// updatePreferences обновляет предпочтения пользователя
func (h *ProfileHandler) updatePreferences(w http.ResponseWriter, userID uuid.UUID, req service.UpdateProfileRequest) {
	updateData := domain.PreferencesUpdateRequest{
		ShowGender:   req.ShowGender,
		AgeMin:       req.AgeMin,
		AgeMax:       req.AgeMax,
		MaxDistance:  req.MaxDistance,
		GlobalSearch: req.GlobalSearch,
	}

	if err := h.profileService.UpdatePreferences(userID, &updateData); err != nil {
		h.logger.Warnf("updatePreferences: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, service.SuccessResponse{Message: "Preferences updated successfully"})
}

// uploadPhotos загружает фотографии
func (h *ProfileHandler) uploadPhotos(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		h.logger.Warnf("uploadPhotos: no files uploaded")
		utils.WriteJSONError(w, http.StatusBadRequest, "no photos provided")
		return
	}

	uploadedPhotos, err := h.profileService.UploadPhotos(userID, files)
	if err != nil {
		status := http.StatusBadRequest
		if err == expectation.ErrPhotoLimitExceeded {
			status = http.StatusForbidden
		}
		h.logger.Warnf("uploadPhotos: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	response := service.UploadPhotosResponse{
		Photos: convertToInterfaceSlice(uploadedPhotos),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// deletePhoto удаляет фотографию
func (h *ProfileHandler) deletePhoto(w http.ResponseWriter, userID uuid.UUID, req service.UpdateProfileRequest) {
	if req.PhotoID == uuid.Nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "photo ID is required")
		return
	}

	if err := h.profileService.DeletePhoto(userID, req.PhotoID); err != nil {
		status := http.StatusBadRequest
		if err == expectation.ErrPhotoNotFound {
			status = http.StatusNotFound
		} else if err == expectation.ErrPhotoNotOwned {
			status = http.StatusForbidden
		}
		h.logger.Warnf("deletePhoto: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, service.SuccessResponse{Message: "Photo deleted successfully"})
}

// setPrimaryPhoto устанавливает основное фото
func (h *ProfileHandler) setPrimaryPhoto(w http.ResponseWriter, userID uuid.UUID, req service.UpdateProfileRequest) {
	if req.PhotoID == uuid.Nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "photo ID is required")
		return
	}

	if err := h.profileService.SetPrimaryPhoto(userID, req.PhotoID); err != nil {
		status := http.StatusBadRequest
		if err == expectation.ErrPhotoNotFound {
			status = http.StatusNotFound
		} else if err == expectation.ErrPhotoNotOwned {
			status = http.StatusForbidden
		}
		h.logger.Warnf("setPrimaryPhoto: %v", err)
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, service.SuccessResponse{Message: "Primary photo set successfully"})
}

// getUserIDFromContext извлекает userID из контекста (будет установлен в middleware)
func (h *ProfileHandler) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	return middleware.GetUserIDFromContext(r.Context())
}

// convertToInterfaceSlice преобразует слайс UserPhoto в interface{}
func convertToInterfaceSlice(photos []domain.UserPhoto) []interface{} {
	result := make([]interface{}, len(photos))
	for i, photo := range photos {
		result[i] = photo
	}
	return result
}
