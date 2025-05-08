package repository

import (
	"context"
	"fmt"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type EmploymentTypeRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewEmploymentTypeRepository(postgres *db.PostgreSQL) EmploymentTypeRepository {
	return &EmploymentTypeRepositoryImpl{
		postgres: postgres,
	}
}

func (r *EmploymentTypeRepositoryImpl) GetAll(ctx context.Context) ([]models.EmploymentType, error) {
	query := `
		SELECT id, name, description
		FROM employment_types
		ORDER BY name
	`

	var employmentTypes []models.EmploymentType
	err := r.postgres.SelectContext(ctx, &employmentTypes, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типов занятости: %w", err)
	}

	return employmentTypes, nil
}

func (r *EmploymentTypeRepositoryImpl) GetByID(ctx context.Context, id int) (*models.EmploymentType, error) {
	query := `
		SELECT id, name, description
		FROM employment_types
		WHERE id = $1
	`

	var employmentType models.EmploymentType
	err := r.postgres.GetContext(ctx, &employmentType, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типа занятости по ID: %w", err)
	}

	return &employmentType, nil
}

func (r *EmploymentTypeRepositoryImpl) GetByName(ctx context.Context, name string) (*models.EmploymentType, error) {
	query := `
		SELECT id, name, description
		FROM employment_types
		WHERE name = $1
	`

	var employmentType models.EmploymentType
	err := r.postgres.GetContext(ctx, &employmentType, query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типа занятости по имени: %w", err)
	}

	return &employmentType, nil
}
