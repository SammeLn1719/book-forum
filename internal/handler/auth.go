package handler

import (
    "net/http"
    "book-forum/internal/service"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Обработка регистрации
    // ...
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // Обработка входа
    // ...
}
