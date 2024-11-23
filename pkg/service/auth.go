package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
)

const salt = "aihvsop198kgmlk"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
