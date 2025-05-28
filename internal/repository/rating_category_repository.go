package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type RatingCategoryRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewRatingCategoryRepository(postgres *db.PostgreSQL) RatingCategoryRepository {
	return &RatingCategoryRepositoryImpl{
		postgres: postgres,
	}
}

func (r *RatingCategoryRepositoryImpl) GetAll(ctx context.Context) ([]models.RatingCategory, error) {
	query := `
		SELECT id, name, description
		FROM rating_categories
		ORDER BY name
	`

	var categories []models.RatingCategory
	err := r.postgres.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении категорий рейтингов: %w", err)
	}

	return categories, nil
}

func (r *RatingCategoryRepositoryImpl) GetByID(ctx context.Context, id int) (*models.RatingCategory, error) {
	query := `
		SELECT id, name, description
		FROM rating_categories
		WHERE id = $1
	`

	var category models.RatingCategory
	err := r.postgres.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении категории рейтинга по ID: %w", err)
	}

	return &category, nil
}

func (r *RatingCategoryRepositoryImpl) GetByName(ctx context.Context, name string) (*models.RatingCategory, error) {
	var category models.RatingCategory
	query := `
		SELECT id, name, created_at, updated_at
		FROM rating_categories
		WHERE name = $1
	`

	err := r.postgres.GetContext(ctx, &category, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("категория не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении категории: %w", err)
	}

	return &category, nil
}

func (r *RatingCategoryRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM rating_categories"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете категорий рейтинга: %w", err)
	}
	return count, nil
}

func (r *RatingCategoryRepositoryImpl) Create(ctx context.Context, category *models.RatingCategory) (int, error) {
	query := `
		INSERT INTO rating_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Description,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании категории рейтинга: %w", err)
	}

	return id, nil
}

func (r *RatingCategoryRepositoryImpl) Update(ctx context.Context, category *models.RatingCategory) error {
	query := `
		UPDATE rating_categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		category.Name,
		category.Description,
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении категории рейтинга: %w", err)
	}

	return nil
}

func (r *RatingCategoryRepositoryImpl) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := `
		SELECT COUNT(*) FROM review_category_ratings WHERE category_id = $1
		UNION ALL
		SELECT COUNT(*) FROM company_category_ratings WHERE category_id = $1
	`

	rows, err := r.postgres.QueryContext(ctx, checkQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при проверке использования категории: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return fmt.Errorf("ошибка при чтении результатов проверки: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("категория используется в рейтингах и не может быть удалена")
		}
	}

	query := "DELETE FROM rating_categories WHERE id = $1"
	_, err = r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении категории рейтинга: %w", err)
	}

	return nil
}
