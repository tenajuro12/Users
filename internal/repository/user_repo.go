package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"users/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.UserCreate) (model.User, error)
	GetByID(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, id string, user model.UserUpdate) (model.User, error)
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, userCreate model.UserCreate) (model.User, error) {
	user := model.User{
		ID:        uuid.New().String(),
		FirstName: userCreate.FirstName,
		LastName:  userCreate.LastName,
		Email:     userCreate.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	query := `INSERT INTO users (id, first_name, last_name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := u.db.Exec(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}
