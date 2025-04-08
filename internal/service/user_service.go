package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"users/internal/model"
	"users/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.UserCreate) (model.User, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, id string, user model.UserUpdate) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}
func (s userService) CreateUser(ctx context.Context, userCreate model.UserCreate) (model.User, error) {
	if err := s.validator.Struct(userCreate); err != nil {
		return model.User{}, err
	}
	return s.userRepo.Create(ctx, userCreate)
}

func (s userService) GetUser(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("id is required")
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s userService) UpdateUser(ctx context.Context, id string, userUpdate model.UserUpdate) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("id is required")
	}

	if err := s.validator.Struct(userUpdate); err != nil {
		return model.User{}, err
	}

	return s.userRepo.Update(ctx, id, userUpdate)
}
func (s userService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return s.userRepo.Delete(ctx, id)
}
