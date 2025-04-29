package usecase

import (
	"book-forum/internal/domain"
	"book-forum/internal/repository"
)

type BookUseCase struct {
	repo *repository.BookRepository
}

func NewBookUseCase(repo *repository.BookRepository) *BookUseCase {
	return &BookUseCase{repo: repo}
}

func (uc *BookUseCase) AddBook(book domain.Book) (int, error) {
	return uc.repo.Create(book)
}

func (uc *BookUseCase) GetBooks() ([]domain.Book, error) {
	return uc.repo.GetAll()
}

func (uc *BookUseCase) InitDB() error {
	return uc.repo.InitTable()
}
