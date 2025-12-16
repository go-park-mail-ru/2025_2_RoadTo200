package utils

import (
	"encoding/json"
	"net/http"

	"github.com/mailru/easyjson"
)

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	var data []byte
	var err error

	// Проверяем, реализует ли тип интерфейс easyjson.Marshaler
	if marshaler, ok := payload.(easyjson.Marshaler); ok {
		data, err = easyjson.Marshal(marshaler)
	} else {
		// Используем стандартный json для типов без easyjson
		data, err = json.Marshal(payload)
	}

	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func ReadJSON(r *http.Request, dst interface{}) error {
	// Проверяем, реализует ли тип интерфейс easyjson.Unmarshaler
	// easyjson.Unmarshaler должен быть указателем
	if unmarshaler, ok := dst.(easyjson.Unmarshaler); ok {
		return easyjson.UnmarshalFromReader(r.Body, unmarshaler)
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func SetSessionCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Message string `json:"message"`
}
