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

func (u *userRepository) Update(ctx context.Context, id string, userUpdate model.UserUpdate) (model.User, error) {
	user, err := u.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	if userUpdate.FirstName != "" {
		user.FirstName = userUpdate.FirstName
	}
	if userUpdate.LastName != "" {
		user.LastName = userUpdate.LastName
	}
	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, updated_at = $4
		WHERE id = $5
	`

	_, err = u.db.Exec(ctx, query,
		user.FirstName, user.LastName, user.Email, user.UpdatedAt, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	commandTag, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
