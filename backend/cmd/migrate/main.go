package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
)

func main() {
	// Замените на свои данные подключения
	connStr := "postgres://admin:12345@localhost:5431/Tinder?sslmode=disable"
	fileToRun := "migrations/dml.sql" // Путь к вашему SQL-файлу

	// Подключение к базе данных
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer conn.Close(ctx)

	// Чтение SQL-файла
	sqlBytes, err := os.ReadFile(fileToRun)
	if err != nil {
		log.Fatalf("Не удалось прочитать SQL-файл: %v", err)
	}

	sqlString := string(sqlBytes)

	// Разделение SQL-файла на отдельные инструкции (на основе точки с запятой, игнорируя `GO` в случае с SQL Server)
	// Простейший подход, но может потребовать более сложной логики для реальных сценариев
	queries := strings.Split(sqlString, ";")

	// Выполнение каждой инструкции
	for _, query := range queries {
		// Удаление пробельных символов в начале и конце инструкции
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := conn.Exec(ctx, query)
		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v\nЗапрос: %s", err, query)
			// Можно прервать выполнение или продолжить
			// return
		} else {
			fmt.Printf("Запрос успешно выполнен: %s\n", query)
		}
	}
}
