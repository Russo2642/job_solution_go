package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"job_solition/internal/db"
	"job_solition/internal/models"
)

type UserRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewUserRepository(postgres *db.PostgreSQL) UserRepository {
	return &UserRepositoryImpl{
		postgres: postgres,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) (int, error) {
	query := `
		INSERT INTO users 
		(email, phone, password_hash, first_name, last_name, role, created_at, updated_at)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	return id, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.postgres.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash, first_name, last_name, role, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := r.postgres.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET phone = $1, first_name = $2, last_name = $3, role = $4, password_hash = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Role,
		user.PasswordHash,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
