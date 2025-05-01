package handler

import (
	"net/http"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"authar"`
	Description string `json:"description"`
	Prise       int    `json:"prise"`
}

func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GetBooks"))
}
