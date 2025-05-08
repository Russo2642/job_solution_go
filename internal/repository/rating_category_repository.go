package repository

import (
	"context"
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
	query := `
		SELECT id, name, description
		FROM rating_categories
		WHERE name = $1
	`

	var category models.RatingCategory
	err := r.postgres.GetContext(ctx, &category, query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении категории рейтинга по имени: %w", err)
	}

	return &category, nil
}
