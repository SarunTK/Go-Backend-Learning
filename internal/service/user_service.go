package service

import (
	"context"
	"fmt"
	"strings"

	"go101/internal/model"
	"go101/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*model.User, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("user repository is not initialized")
	}

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("user repository is not initialized")
	}
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("user repository is not initialized")
	}
	return s.repo.List(ctx)
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
