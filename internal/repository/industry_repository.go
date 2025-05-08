package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"job_solition/internal/db"
	"job_solition/internal/models"
)

type IndustryRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewIndustryRepository(postgres *db.PostgreSQL) IndustryRepository {
	return &IndustryRepositoryImpl{
		postgres: postgres,
	}
}

func (r *IndustryRepositoryImpl) GetAll(ctx context.Context, filter models.IndustryFilter) ([]models.Industry, int, error) {
	baseQuery := `
		FROM industries 
		WHERE 1=1
	`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	queryConditions := baseQuery
	if len(conditions) > 0 {
		queryConditions += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + queryConditions

	sortBy := "name"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	sortOrder := "ASC"
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	dataQuery := fmt.Sprintf(`
		SELECT id, name, color
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете отраслей: %w", err)
	}

	var industries []models.Industry
	err = r.postgres.SelectContext(ctx, &industries, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении отраслей: %w", err)
	}

	return industries, total, nil
}

func (r *IndustryRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Industry, error) {
	query := `
		SELECT id, name, color
		FROM industries 
		WHERE id = $1
	`

	var industry models.Industry
	err := r.postgres.GetContext(ctx, &industry, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("отрасль не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении отрасли: %w", err)
	}

	return &industry, nil
}

func (r *IndustryRepositoryImpl) GetByIDs(ctx context.Context, ids []int) ([]models.Industry, error) {
	if len(ids) == 0 {
		return []models.Industry{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, name, color
		FROM industries 
		WHERE id IN (%s)
		ORDER BY name
	`, strings.Join(placeholders, ", "))

	var industries []models.Industry
	err := r.postgres.SelectContext(ctx, &industries, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отраслей: %w", err)
	}

	return industries, nil
}

func (r *IndustryRepositoryImpl) GetByCompanyID(ctx context.Context, companyID int) ([]models.Industry, error) {
	query := `
		SELECT i.id, i.name, i.color
		FROM industries i
		JOIN company_industries ci ON i.id = ci.industry_id
		WHERE ci.company_id = $1
		ORDER BY i.name
	`

	var industries []models.Industry
	err := r.postgres.SelectContext(ctx, &industries, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отраслей компании: %w", err)
	}

	return industries, nil
}

func (r *IndustryRepositoryImpl) AddCompanyIndustry(ctx context.Context, companyID, industryID int) error {
	query := `
		INSERT INTO company_industries (company_id, industry_id)
		VALUES ($1, $2)
		ON CONFLICT (company_id, industry_id) DO NOTHING
	`

	_, err := r.postgres.ExecContext(ctx, query, companyID, industryID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении отрасли компании: %w", err)
	}

	return nil
}

func (r *IndustryRepositoryImpl) RemoveCompanyIndustry(ctx context.Context, companyID, industryID int) error {
	query := `
		DELETE FROM company_industries
		WHERE company_id = $1 AND industry_id = $2
	`

	_, err := r.postgres.ExecContext(ctx, query, companyID, industryID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении отрасли компании: %w", err)
	}

	return nil
}

func (r *IndustryRepositoryImpl) UpdateColor(ctx context.Context, id int, color string) error {
	query := `
		UPDATE industries
		SET color = $1
		WHERE id = $2
	`

	result, err := r.postgres.ExecContext(ctx, query, color, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении цвета отрасли: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("отрасль не найдена")
	}

	return nil
}
