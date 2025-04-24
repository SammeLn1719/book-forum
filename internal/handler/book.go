package handler

import (
    "net/http"
    "book-forum/internal/service"
)

type BookHandler struct {
    bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
    return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    // Реализация обработчика
}
