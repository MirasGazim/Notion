package service

import "notion/internal/repository"

type Service struct {
	repository.Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
