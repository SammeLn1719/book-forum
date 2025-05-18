package handler

import (
	"encoding/json"
	"net/http"
	"time"
    "strings"
	"github.com/google/uuid"
	"book-forum/internal/models"
	"book-forum/internal/repository"
	"book-forum/internal/utils"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// ValidateRegistrationRequest проверяет все обязательные поля
func ValidateRegistrationRequest(req RegisterRequest) []ValidationError {
    var errors []ValidationError

    // Проверка каждого поля с кастомными сообщениями
    if strings.TrimSpace(req.Username) == "" {
        errors = append(errors, ValidationError{
            Field:   "username",
            Message: "Имя пользователя обязательно для заполнения",
        })
    }

    if strings.TrimSpace(req.Email) == "" {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: "Email обязателен для заполнения",
        })
    }

    if strings.TrimSpace(req.PasswordHash) == "" {
        errors = append(errors, ValidationError{
            Field:   "password_hash",
            Message: "Пароль обязателен для заполнения",
        })
    }

    // Дополнительные проверки (например, формат email)
    if req.Email != "" && !strings.Contains(req.Email, "@") {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: "Некорректный формат email",
        })
    }

    return errors
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    // Валидация
// Валидация
    if validationErrors := ValidateRegistrationRequest(req); len(validationErrors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "errors": validationErrors,
        })
        return
    }

    if len(user.PasswordHash) < 8 {
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
    user.PasswordHash = ""

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        PasswordHash string `json:"password_hash"`
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

    if !utils.CheckPasswordHash(creds.PasswordHash, user.PasswordHash) {
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
