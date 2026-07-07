package service

import (
	"errors"
	"notion/internal/repository"
)

type Service struct {
	repository.Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)
