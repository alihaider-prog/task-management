package service

import (
	"errors"
	"task-management/internal/models"
	"task-management/internal/repository"
	"task-management/pkg/jwt"
	"task-management/pkg/password"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) error {
	hash, err := password.Hash(req.Passowrd)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
	}

	return s.repo.Create(user)
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponce, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := password.Check(req.Passowrd, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponce{
		Token: token,
	}, nil
}
