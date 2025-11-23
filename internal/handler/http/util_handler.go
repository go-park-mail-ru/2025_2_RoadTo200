package handler

import (
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

// HealthHandler godoc
// @Summary Проверка здоровья приложения
// @Description Возвращает текущий статус работы приложения Terabithia. Используется для мониторинга и проверки доступности сервиса.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Приложение работает корректно"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/health [get]
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"app":    "terabithia app",
	})
}

func RegisterSwagger(url string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/doc.json" {
			data, err := os.ReadFile(url)
			if err != nil {
				logger.Error("Error reading swagger.json: %v", err)
				utils.WriteJSONError(w, http.StatusInternalServerError, "Swagger docs not found: "+err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")).ServeHTTP(w, r)
	}
}
