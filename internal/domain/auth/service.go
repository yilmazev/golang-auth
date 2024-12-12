package auth

import (
	"errors"
	"golang-auth/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *AuthRepository
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTSecret())
}

func (s *AuthService) Register(user User) (string, error) {
	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	if err := s.Repo.CreateUser(user); err != nil {
		return "", err
	}

	return s.GenerateToken(user.Username)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := s.CheckPasswordHash(password, user.Password); err != nil {
		return "", errors.New("invalid username or password")
	}

	return s.GenerateToken(username)
}

func (s *AuthService) GetUserByUsername(username string) (*User, error) {
	return s.Repo.GetUserByUsername(username)
}

func (s *AuthService) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return username, nil
	}
	return "", errors.New("invalid token")
}
