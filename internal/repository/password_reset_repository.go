package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"job_solition/internal/db"
	"job_solition/internal/models"
)

type PasswordResetRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewPasswordResetRepository(postgres *db.PostgreSQL) PasswordResetRepository {
	return &PasswordResetRepositoryImpl{
		postgres: postgres,
	}
}

func (r *PasswordResetRepositoryImpl) Create(ctx context.Context, token *models.PasswordResetToken) (int, error) {
	if err := r.DeleteByUserID(ctx, token.UserID); err != nil {
		return 0, fmt.Errorf("ошибка при удалении существующих токенов: %w", err)
	}

	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.postgres.GetContext(
		ctx,
		&id,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании токена сброса пароля: %w", err)
	}

	return id, nil
}

func (r *PasswordResetRepositoryImpl) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	var resetToken models.PasswordResetToken
	err := r.postgres.GetContext(ctx, &resetToken, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("токен сброса пароля не найден")
		}
		return nil, fmt.Errorf("ошибка при получении токена сброса пароля: %w", err)
	}

	return &resetToken, nil
}

func (r *PasswordResetRepositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE token = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении токена сброса пароля: %w", err)
	}

	return nil
}

func (r *PasswordResetRepositoryImpl) DeleteByUserID(ctx context.Context, userID int) error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE user_id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении токенов сброса пароля пользователя: %w", err)
	}

	return nil
}
