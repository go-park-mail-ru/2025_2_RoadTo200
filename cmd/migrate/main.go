package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/jackc/pgx/v4"
)

func main() {
	logger.Printf("Args amount %d", len(os.Args)-1)
	// Замените на свои данные подключения
	connStr := "postgres://admin:12345@localhost:5431/Tinder?sslmode=disable"
	fileToRun := "migrations/%s.sql"
	if len(os.Args) == 2 && os.Args[1] != "" {
		fileToRun = fmt.Sprintf(fileToRun, os.Args[1])
	} else {
		logger.Fatal("Enter arg")
	}

	// Подключение к базе данных
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		logger.Fatal("Не удалось подключиться к базе данных: %v", err)
	}
	defer conn.Close(ctx)

	// Чтение SQL-файла
	sqlBytes, err := os.ReadFile(fileToRun)
	if err != nil {
		logger.Fatal("Не удалось прочитать SQL-файл: %v", err)
	}

	sqlString := string(sqlBytes)

	// Разделение SQL-файла на отдельные инструкции (на основе точки с запятой, игнорируя `GO` в случае с SQL Server)
	// Простейший подход, но может потребовать более сложной логики для реальных сценариев
	queries := strings.Split(sqlString, ";\n")

	buf := ""
	mod := 'r'
	// Выполнение каждой инструкции
	for _, query := range queries {
		// Удаление пробельных символов в начале и конце инструкции
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if mod == 'f' || strings.HasPrefix(query, "-- Функция") {
			mod = 'f'
			if query != "-- end" {
				buf += query + ";"
				continue
			} else {
				query = buf
				mod = 'r'
			}
		}

		_, err := conn.Exec(ctx, query)
		if err != nil {
			logger.Error("Ошибка при выполнении запроса: %v\nЗапрос: %s", err, query)
			// Можно прервать выполнение или продолжить
			// return
		} else {
			logger.Printf("Запрос успешно выполнен: %s", query)
		}
	}
}
