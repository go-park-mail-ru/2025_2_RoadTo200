package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	//"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile возвращает профиль текущего пользователя
// @Summary Get user profile
// @Description Get current user's profile with photos and preferences
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} service.ProfileResponse
// @Failure 401 {object} service.ErrorResponse
// @Failure 500 {object} service.ErrorResponse
// @Router /profile/profile [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		if err == errors.ErrProfileNotFound {
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

// ChangeProfile обрабатывает все операции изменения профиля
// @Summary Update profile
// @Description Update user profile information, upload photos, delete photos, set primary photo
// @Tags profile
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param action formData string true "Action type: updateInfo, uploadPhotos, deletePhoto, setPrimaryPhoto, updatePreferences"
// @Param name formData string false "User name"
// @Param phone formData string false "Phone number"
// @Param birth_date formData string false "Birth date (YYYY-MM-DD)"
// @Param gender formData string false "Gender"
// @Param bio formData string false "Bio"
// @Param latitude formData number false "Latitude"
// @Param longitude formData number false "Longitude"
// @Param show_gender formData string false "Gender preference"
// @Param age_min formData int false "Minimum age preference"
// @Param age_max formData int false "Maximum age preference"
// @Param max_distance formData int false "Maximum distance preference"
// @Param global_search formData bool false "Global search preference"
// @Param photo_id formData string false "Photo ID for delete/set primary"
// @Param photos formData file false "Photos to upload"
// @Success 200 {object} interface{}
// @Failure 400 {object} service.ErrorResponse
// @Failure 401 {object} service.ErrorResponse
// @Failure 500 {object} service.ErrorResponse
// @Router /profile/changeProfile [post]
func (h *ProfileHandler) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))
	fmt.Printf("Method: %s\n", r.Method)
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	contentType := r.Header.Get("Content-Type")
	fmt.Printf("Content-Type: %s\n", contentType)

	// Используем strings.Contains вместо строгого сравнения
	if strings.Contains(contentType, "application/json") {
		h.handleJSONRequest(w, r, userID)
	} else if strings.Contains(contentType, "multipart/form-data") {
		h.handleMultipartRequest(w, r, userID)
	} else {
		utils.WriteJSONError(w, http.StatusBadRequest, "unsupported content type: "+contentType)
	}
}

// handleJSONRequest обрабатывает JSON запросы
func (h *ProfileHandler) handleJSONRequest(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var req service.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

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
		utils.WriteJSONError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	action := r.FormValue("action")
	if action != "uploadPhotos" {
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
		if err == errors.ErrProfileNotFound {
			status = http.StatusNotFound
		}
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
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, service.SuccessResponse{Message: "Preferences updated successfully"})
}

// uploadPhotos загружает фотографии
func (h *ProfileHandler) uploadPhotos(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "no photos provided")
		return
	}

	uploadedPhotos, err := h.profileService.UploadPhotos(userID, files)
	if err != nil {
		status := http.StatusBadRequest
		if err == errors.ErrPhotoLimitExceeded {
			status = http.StatusForbidden
		}
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
		if err == errors.ErrPhotoNotFound {
			status = http.StatusNotFound
		} else if err == errors.ErrPhotoNotOwned {
			status = http.StatusForbidden
		}
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
		if err == errors.ErrPhotoNotFound {
			status = http.StatusNotFound
		} else if err == errors.ErrPhotoNotOwned {
			status = http.StatusForbidden
		}
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
