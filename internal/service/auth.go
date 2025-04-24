package service

import (
    "book-forum/internal/models"
    "book-forum/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *AuthService) Register(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.PasswordHash = string(hashedPassword)
    return s.userRepo.CreateUser(context.Background(), user)
}
