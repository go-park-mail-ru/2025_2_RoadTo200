package main

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/app"
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

// @host 217.16.17.116:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey SessionToken
// @in header
// @name X-Session-Token
// @description Токен сессии для аутентификации пользователя (альтернатива cookie)

func main() {
	//from := "ender-fox@mail.ru"
	//password := "your_password" // Пароль приложения (для двухфакторной аутентификации)
	//to := []string{"your_password"}
	//smtpHost := "smtp.mail.ru"
	//smtpPort := "587"
	//
	//// Заголовки письма
	//subject := "Subject: Тестовое письмо из Go!\n"
	//body := "Текст письма.\n"
	//msg := []byte(subject + "\n" + body)
	//
	//// Настройка TLS
	//tlsConfig := &tls.Config{
	//	ServerName: smtpHost,
	//}
	//
	//// Подключение к серверу
	//conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//
	//client, err := smtp.NewClient(conn, smtpHost)
	//if err != nil {
	//	panic(err)
	//}
	//defer client.Close()
	//
	//// Аутентификация
	//auth := smtp.PlainAuth("", from, password, smtpHost)
	//if err = client.Auth(auth); err != nil {
	//	panic(err)
	//}
	//
	//// Отправка
	//if err = client.Mail(from); err != nil {
	//	panic(err)
	//}
	//for _, addr := range to {
	//	if err = client.Rcpt(addr); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//w, err := client.Data()
	//if err != nil {
	//	panic(err)
	//}
	//_, err = w.Write(msg)
	//if err != nil {
	//	panic(err)
	//}
	//err = w.Close()
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Письмо отправлено!")
	app.Run()
}
