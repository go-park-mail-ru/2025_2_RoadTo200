// вспомогательные функции (JSON-ответы и т.п.)
package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON - это вспомогательная функция для записи ответа в формате JSON
// Она устанавливает заголовок Content-Type и записывает код состояния и полезную нагрузку
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	var data []byte
	var err error
	switch payload.(type) {
	case []byte:
		data = payload.([]byte)
	default:
		data, err = json.Marshal(payload)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Write(data)
}

// readJSON - это вспомогательная функция для чтения текста запроса в формате JSON
// Она декодирует текст в формате JSON в предоставленный целевой интерфейс
func readJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
