package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"book-forum/internal/models"
	"book-forum/internal/repository"
	"book-forum/internal/utils"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    // Валидация
    if user.Username == "" || user.Email == "" || user.Password == "" {
        http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
        return
    }

    if len(user.Password) < 8 {
        http.Error(w, "Пароль должен содержать минимум 8 символов", http.StatusBadRequest)
        return
    }

    // Проверка уникальности email
    exists, err := h.userRepo.IsEmailExists(user.Email)
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        return
    }
    if exists {
        http.Error(w, "Email уже зарегистрирован", http.StatusConflict)
        return
    }

    // Создание пользователя
    if err := h.userRepo.CreateUser(&user); err != nil {
        http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
        return
    }

    // Убираем пароль из ответа
    user.Password = ""

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.userRepo.GetUserByEmail(creds.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if !utils.CheckPasswordHash(creds.Password, user.Password) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    session := models.Session{
        ID:        uuid.New().String(),
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    if err := h.sessionRepo.CreateSession(&session); err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    session.ID,
        Expires:  session.ExpiresAt,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "Already logged out", http.StatusBadRequest)
        return
    }

    if err := h.sessionRepo.DeleteSession(cookie.Value); err != nil {
        http.Error(w, "Failed to logout", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
    })

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
