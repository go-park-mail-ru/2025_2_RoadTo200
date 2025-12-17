package main

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/app_run/gateway-service/app"
)

//тест

// @title Terabithia Dating App API
// @version 1.0
// @description API для dating приложения Terabithia
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@terabithia.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey SessionToken
// @in header
// @name X-Session-Token
// @description Токен сессии для аутентификации пользователя (альтернатива cookie)

func main() {
	app.Run()
}
