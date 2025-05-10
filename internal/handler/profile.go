package handler

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	// Здесь реализуйте получение данных профиля из базы
	w.Write([]byte("Profile page for user ID: " + string(userID)))
}
