package repository

import (
	"context"
	"fmt"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type EmploymentPeriodRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewEmploymentPeriodRepository(postgres *db.PostgreSQL) EmploymentPeriodRepository {
	return &EmploymentPeriodRepositoryImpl{
		postgres: postgres,
	}
}

func (r *EmploymentPeriodRepositoryImpl) GetAll(ctx context.Context) ([]models.EmploymentPeriod, error) {
	query := `
		SELECT id, name, description
		FROM employment_periods
		ORDER BY name
	`

	var employmentPeriods []models.EmploymentPeriod
	err := r.postgres.SelectContext(ctx, &employmentPeriods, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении периодов работы: %w", err)
	}

	return employmentPeriods, nil
}

func (r *EmploymentPeriodRepositoryImpl) GetByID(ctx context.Context, id int) (*models.EmploymentPeriod, error) {
	query := `
		SELECT id, name, description
		FROM employment_periods
		WHERE id = $1
	`

	var employmentPeriod models.EmploymentPeriod
	err := r.postgres.GetContext(ctx, &employmentPeriod, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении периода работы по ID: %w", err)
	}

	return &employmentPeriod, nil
}

func (r *EmploymentPeriodRepositoryImpl) GetByName(ctx context.Context, name string) (*models.EmploymentPeriod, error) {
	query := `
		SELECT id, name, description
		FROM employment_periods
		WHERE name = $1
	`

	var employmentPeriod models.EmploymentPeriod
	err := r.postgres.GetContext(ctx, &employmentPeriod, query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении периода работы по имени: %w", err)
	}

	return &employmentPeriod, nil
}
