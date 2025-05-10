package middleware

import (
	"context"
	"net/http"
	"time"
	"book-forum/internal/repository"
)

func AuthMiddleware(sessionRepo *repository.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := sessionRepo.GetSession(cookie.Value)
			if err != nil || session.ExpiresAt.Before(time.Now()) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
