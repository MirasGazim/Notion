package service

import (
	"context"
	"errors"
	"fmt"
	"notion/internal/models/user"
	"notion/internal/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "dfhgsdfhgidu1224"
	signingKey = "grkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, u user.SignUpRequest) (uuid.UUID, error) {
	const op = "service/auth/CreateUser"
	hash, err := generatePasswordHash(u.Password)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	u.Password = hash
	return s.repo.CreateUser(ctx, u)
}

func (s *AuthService) GetUser(ctx context.Context, u user.SignInRequest) (user.AuthUser, error) {
	const op = "service/auth/GetUser"
	id, err := s.repo.GetUser(ctx, u)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return user.AuthUser{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return user.AuthUser{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(id.Password), []byte(u.Password))
	if err != nil {
		return user.AuthUser{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	return id, nil
}

// func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
// 	u, err := s.repo.GetUser(ctx, username)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
// 	if err != nil {
// 		return "", errors.New("invalid username or password")
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 		u.ID,
// 	})

// 	return token.SignedString([]byte(signingKey))
// }

// func (s *AuthService) ParseToken(ctx context.Context, accesstoken string) (uuid.UUID, error) {
// 	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid signing method")
// 		}
// 		return []byte(signingKey), nil
// 	})

// 	if err != nil {
// 		return uuid.UUID{}, err
// 	}

// 	claims, ok := token.Claims.(*tokenClaims)
// 	if !ok {
// 		return uuid.UUID{}, errors.New("token claims are not of type *tokenClaims")
// 	}
// 	return claims.UserID, nil
// }

func generatePasswordHash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}
