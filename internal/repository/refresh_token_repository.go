package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"job_solition/internal/db"
	"job_solition/internal/models"
)

type RefreshTokenRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewRefreshTokenRepository(postgres *db.PostgreSQL) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{
		postgres: postgres,
	}
}

func (r *RefreshTokenRepositoryImpl) Create(ctx context.Context, token *models.RefreshToken) (int, error) {
	query := `
		INSERT INTO refresh_tokens 
		(user_id, token, expires_at, created_at)
		VALUES 
		($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании refresh токена: %w", err)
	}

	return id, nil
}

func (r *RefreshTokenRepositoryImpl) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens 
		WHERE token = $1
	`

	var refreshToken models.RefreshToken
	err := r.postgres.GetContext(ctx, &refreshToken, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh токен не найден")
		}
		return nil, fmt.Errorf("ошибка при получении refresh токена: %w", err)
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении refresh токена: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID int) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении refresh токенов пользователя: %w", err)
	}

	return nil
}
