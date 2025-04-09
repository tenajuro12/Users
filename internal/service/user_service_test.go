package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"users/internal/model"
	"users/internal/service"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user model.UserCreate) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, user model.UserUpdate) (model.User, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userCreate := model.UserCreate{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}

		expectedUser := model.User{
			ID:        uuid.New().String(),
			FirstName: userCreate.FirstName,
			LastName:  userCreate.LastName,
			Email:     userCreate.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Create", ctx, userCreate).Return(expectedUser, nil).Once()

		user, err := userService.CreateUser(ctx, userCreate)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		invalidUser := model.UserCreate{
			FirstName: "John",
			// Нет фамилии
			Email: "john.doe@example.com",
		}

		user, err := userService.CreateUser(ctx, invalidUser)

		assert.Error(t, err)
		assert.Empty(t, user)
	})

	t.Run("Repository error", func(t *testing.T) {
		userCreate := model.UserCreate{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}

		expectedErr := errors.New("database error")
		mockRepo.On("Create", ctx, userCreate).Return(model.User{}, expectedErr).Once()

		user, err := userService.CreateUser(ctx, userCreate)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		expectedUser := model.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil).Once()

		user, err := userService.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		user, err := userService.GetUser(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, "id is required", err.Error())
		assert.Empty(t, user)
	})

	t.Run("Repository error", func(t *testing.T) {
		userID := uuid.New().String()
		expectedErr := errors.New("user not found")

		mockRepo.On("GetByID", ctx, userID).Return(model.User{}, expectedErr).Once()

		user, err := userService.GetUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		userUpdate := model.UserUpdate{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
		}

		expectedUser := model.User{
			ID:        userID,
			FirstName: userUpdate.FirstName,
			LastName:  userUpdate.LastName,
			Email:     userUpdate.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Update", ctx, userID, userUpdate).Return(expectedUser, nil).Once()

		user, err := userService.UpdateUser(ctx, userID, userUpdate)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		userUpdate := model.UserUpdate{
			FirstName: "Jane",
		}

		user, err := userService.UpdateUser(ctx, "", userUpdate)

		assert.Error(t, err)
		assert.Equal(t, "id is required", err.Error())
		assert.Empty(t, user)
	})

	t.Run("Validation error", func(t *testing.T) {
		userID := uuid.New().String()
		invalidUpdate := model.UserUpdate{
			Email: "invalid-email", // формат емайл
		}

		user, err := userService.UpdateUser(ctx, userID, invalidUpdate)

		assert.Error(t, err)
		assert.Empty(t, user)
	})

	t.Run("Repository error", func(t *testing.T) {
		userID := uuid.New().String()
		userUpdate := model.UserUpdate{
			FirstName: "Jane",
		}

		expectedErr := errors.New("user not found")
		mockRepo.On("Update", ctx, userID, userUpdate).Return(model.User{}, expectedErr).Once()

		user, err := userService.UpdateUser(ctx, userID, userUpdate)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()

		mockRepo.On("Delete", ctx, userID).Return(nil).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		err := userService.DeleteUser(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, "id is required", err.Error())
	})

	t.Run("Repository error", func(t *testing.T) {
		userID := uuid.New().String()
		expectedErr := errors.New("user not found")

		mockRepo.On("Delete", ctx, userID).Return(expectedErr).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
