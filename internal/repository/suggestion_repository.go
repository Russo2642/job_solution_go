package repository

import (
	"context"
	"fmt"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type SuggestionRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewSuggestionRepository(postgres *db.PostgreSQL) SuggestionRepository {
	return &SuggestionRepositoryImpl{
		postgres: postgres,
	}
}

func (r *SuggestionRepositoryImpl) Create(ctx context.Context, suggestion *models.Suggestion) (int, error) {
	query := `
		INSERT INTO suggestions (type, text, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		suggestion.Type,
		suggestion.Text,
		suggestion.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании предложения: %w", err)
	}

	return id, nil
}

func (r *SuggestionRepositoryImpl) GetAll(ctx context.Context, filter models.SuggestionFilter) ([]models.Suggestion, int, error) {
	baseQuery := `FROM suggestions`
	whereClause := ""
	args := []interface{}{}
	argID := 1

	if filter.Type != "" {
		whereClause = fmt.Sprintf(" WHERE type = $%d", argID)
		args = append(args, filter.Type)
		argID++
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.Limit

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s%s", baseQuery, whereClause)
	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении количества предложений: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, type, text, created_at
		%s%s
		ORDER BY created_at %s
		LIMIT $%d OFFSET $%d
	`, baseQuery, whereClause, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var suggestions []models.Suggestion
	err = r.postgres.SelectContext(ctx, &suggestions, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении предложений: %w", err)
	}

	return suggestions, total, nil
}

func (r *SuggestionRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM suggestions WHERE id = $1`
	_, err := r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении предложения: %w", err)
	}
	return nil
}

func (r *SuggestionRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM suggestions"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете предложений: %w", err)
	}
	return count, nil
}
