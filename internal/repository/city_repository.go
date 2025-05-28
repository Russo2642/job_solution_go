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

type CityRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewCityRepository(postgres *db.PostgreSQL) CityRepository {
	return &CityRepositoryImpl{
		postgres: postgres,
	}
}

func (r *CityRepositoryImpl) GetAll(ctx context.Context, filter models.CityFilter) ([]models.City, int, error) {
	baseQuery := `
		FROM cities 
		WHERE 1=1
	`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR region ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	if filter.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", argID))
		args = append(args, filter.Country)
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
		SELECT id, name, region, country
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете городов: %w", err)
	}

	var cities []models.City
	err = r.postgres.SelectContext(ctx, &cities, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении городов: %w", err)
	}

	return cities, total, nil
}

func (r *CityRepositoryImpl) GetByID(ctx context.Context, id int) (*models.City, error) {
	query := `
		SELECT id, name, region, country
		FROM cities 
		WHERE id = $1
	`

	var city models.City
	err := r.postgres.GetContext(ctx, &city, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("город не найден")
		}
		return nil, fmt.Errorf("ошибка при получении города: %w", err)
	}

	return &city, nil
}

func (r *CityRepositoryImpl) Search(ctx context.Context, query string) ([]models.City, error) {
	q := `
		SELECT id, name, created_at, updated_at
		FROM cities
		WHERE LOWER(name) LIKE LOWER($1)
		ORDER BY name
	`

	var cities []models.City
	err := r.postgres.SelectContext(ctx, &cities, q, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске городов: %w", err)
	}

	return cities, nil
}

func (r *CityRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM cities"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете городов: %w", err)
	}
	return count, nil
}

func (r *CityRepositoryImpl) Create(ctx context.Context, city *models.City) (int, error) {
	query := `
		INSERT INTO cities (name, region, country, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		city.Name,
		city.Region,
		city.Country,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании города: %w", err)
	}

	return id, nil
}

func (r *CityRepositoryImpl) Update(ctx context.Context, city *models.City) error {
	query := `
		UPDATE cities
		SET name = $1, region = $2, country = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		city.Name,
		city.Region,
		city.Country,
		city.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении города: %w", err)
	}

	return nil
}

func (r *CityRepositoryImpl) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := `
		SELECT COUNT(*) FROM companies WHERE city_id = $1
		UNION ALL
		SELECT COUNT(*) FROM users WHERE city_id = $1
	`

	rows, err := r.postgres.QueryContext(ctx, checkQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при проверке использования города: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return fmt.Errorf("ошибка при чтении результатов проверки: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("город используется в компаниях или профилях пользователей и не может быть удален")
		}
	}

	query := "DELETE FROM cities WHERE id = $1"
	_, err = r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении города: %w", err)
	}

	return nil
}
