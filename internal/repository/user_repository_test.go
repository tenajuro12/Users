package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"users/internal/model"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	//Local db check
	if dbURL == "" {
		dbURL = "postgres://postgres:murderpe@localhost:5432/user"
	}

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	_, err = dbpool.Exec(context.Background(), "TRUNCATE users")
	require.NoError(t, err)

	return dbpool
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	userCreate := model.UserCreate{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	user, err := repo.Create(context.Background(), userCreate)
	require.NoError(t, err)

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, userCreate.FirstName, user.FirstName)
	assert.Equal(t, userCreate.LastName, user.LastName)
	assert.Equal(t, userCreate.Email, user.Email)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	userCreate := model.UserCreate{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	createdUser, err := repo.Create(context.Background(), userCreate)
	require.NoError(t, err)

	user, err := repo.GetByID(context.Background(), createdUser.ID)
	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, userCreate.FirstName, user.FirstName)
	assert.Equal(t, userCreate.LastName, user.LastName)
	assert.Equal(t, userCreate.Email, user.Email)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	userCreate := model.UserCreate{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	createdUser, err := repo.Create(context.Background(), userCreate)
	require.NoError(t, err)

	userUpdate := model.UserUpdate{
		FirstName: "Jane",
		Email:     "jane.doe@example.com",
	}

	updatedUser, err := repo.Update(context.Background(), createdUser.ID, userUpdate)
	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, userUpdate.FirstName, updatedUser.FirstName)
	assert.Equal(t, userCreate.LastName, updatedUser.LastName) // LastName didn't change
	assert.Equal(t, userUpdate.Email, updatedUser.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	userCreate := model.UserCreate{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	createdUser, err := repo.Create(context.Background(), userCreate)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), createdUser.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), createdUser.ID)
	assert.Error(t, err)
}
