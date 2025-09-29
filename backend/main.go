// запуск сервера + хендлеры

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// registerHandler обрабатывает запросы на регистрацию пользователей
// Он ожидает POST-запрос с адресом электронной почты, паролем и passwordConfirm в теле JSON
// В случае успеха он создает нового пользователя, запускает сеанс и устанавливает сессионный файл cookie
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не разрешен"})
		return
	}

	var req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}

	//проверка, пришел ли вообще json и читаем его
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное тело запроса"})
		return
	}

	//проверка email
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный email"})
		return
	}

	//проверка пароля
	if len(req.Password) < 6 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароль слишком короткий"})
		return
	}
	// проверка что пароль содержит буквы и цифры
	hasLetter := strings.ContainsAny(req.Password, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasDigit := strings.ContainsAny(req.Password, "0123456789")
	if !hasLetter || !hasDigit {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароль должен содержать буквы и цифры"})
		return
	}
	if req.Password != req.PasswordConfirm {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароли не совпадают"})
		return
	}

	//проверка что пользователь существует
	_, exists := findUser(req.Email)
	if exists {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Пользователь уже зарегистрирован"})
		return
	}

	user, err := createUser(req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Внутренняя ошибка"})
		return
	}

	session := createSession(user.Email, 3600*time.Second) // 1 час
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600, // 1 час
	})

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})

	fmt.Printf("%v: Registered user: %s, %s\n", time.Now(), user.Email, user.Password)
}

// loginHandler обрабатывает запросы пользователя на вход в систему
// Он ожидает POST-запрос с адресом электронной почты и паролем в формате JSON
// При успешной аутентификации он запускает сеанс и устанавливает сессионный файл cookie
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не разрешен"})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное тело запроса"})
		return
	}

	user, exists := findUser(req.Email)
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный email или пароль"})
		return
	}

	//проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	session := createSession(user.Email, 3600*time.Second)

	fmt.Printf("%v: User: %s entered with session: %s\n", time.Now(), user.Email, session.Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})

	fmt.Println("Login user:", user.Email, user.Password)
}

// SessionHandler проверяет статус сеанса текущего пользователя
// Он ожидает запрос GET с файлом cookie "session_token"
// Он возвращает, прошел ли пользователь проверку подлинности, и информацию о пользователе, если да
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не разрешен"})
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	session, ok := getSession(cookie.Value)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	user, ok := findUser(session.UserEmail)
	if !ok {
		fmt.Printf("%v: Error: User %s not found\n", time.Now(), session.UserEmail)
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"authenticated": false,
		})
		return
	}
	fmt.Printf("%v: User: %s entered with session: %s\n", time.Now(), user.ID, session.Token)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// LogoutHandler обрабатывает запросы пользователя на выход из системы
// Он ожидает POST-запрос с файлом cookie "session_token"
// Он удаляет сеанс и очищает файл cookie сеансd
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не разрешен"})
		return
	}

	// читаем куку
	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Нет сессии"})
		return
	}

	// удаляем сессию
	deleteSession(cookie.Value)

	fmt.Printf("%v: Session - %s closed.\n", time.Now(), cookie.Value)
	// очищаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // удалить
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "Выход выполнен"})
}

// feedHandler предоставляет доступ к профилям пользователей для прошедшего проверку подлинности пользователя
// Для этого требуется действительный сеанс
func feedHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка сессии (использует твой getSession)
	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
		return
	}
	_, ok := getSession(cookie.Value)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
		return
	}
	fmt.Printf("%v: Session - %s data sending.\n", time.Now(), cookie.Value)
	writeJSON(w, http.StatusOK, []byte(cards))
}

// swipeHandler обрабатывает действия пользователя по переходу в другой профиль
// Для этого требуется действительный сеанс и ожидается запрос POST с идентификатором цели и направлением
func swipeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
		return
	}
	session, ok := getSession(cookie.Value)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Не авторизован"})
		return
	}

	// Считываем свайп
	var req struct {
		TargetID  string `json:"target_id"`
		Direction string `json:"direction"` // "left" или "right"
	}
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный запрос"})
		return
	}

	// Пока просто логируем — без логики
	fmt.Printf("%s swiped %s %s\n", session.UserEmail, req.Direction, req.TargetID)

	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func addHandler(path string, f func(w http.ResponseWriter, r *http.Request)) {
	http.Handle(path, CORSMiddleware(http.HandlerFunc(f)))
}

// main - это точка входа приложения
// Она устанавливает HTTP-маршруты и запускает сервер
func main() {
	addHandler("/api/register", registerHandler)
	addHandler("/api/login", loginHandler)
	addHandler("/api/session", sessionHandler)
	addHandler("/api/logout", logoutHandler)
	addHandler("/api/feed", feedHandler)
	addHandler("/api/swipe", swipeHandler)

	fmt.Println("Server running on http://:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
