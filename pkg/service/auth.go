package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo       repository.Authorization
	signingKey string
	salt       string
}

func NewAuthService(repo repository.Authorization) (*AuthService, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	signningKey := os.Getenv("SIGNING_KEY")
	salt := os.Getenv("SALT")
	return &AuthService{repo: repo, signingKey: signningKey, salt: salt}, nil
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(a.salt)))
}

func (a *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := a.repo.GetUser(username, a.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	tokenStr, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(a.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not if type *tokenClaims")
	}

	return claims.UserId, nil
}
