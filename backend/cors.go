package main

import (
	"net/http"
)

const origin = "http://127.0.0.1:8001"

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin (замените на конкретный домен для продакшена)
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// Разрешаем методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Разрешаем заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обработка предварительного запроса (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
