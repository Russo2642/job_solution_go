package repository

import (
	"context"
	"database/sql"
	"errors"
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
		SELECT id, name, created_at, updated_at
		FROM employment_periods
		WHERE name = $1
	`

	var period models.EmploymentPeriod
	err := r.postgres.GetContext(ctx, &period, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("период не найден")
		}
		return nil, fmt.Errorf("ошибка при получении периода: %w", err)
	}

	return &period, nil
}

func (r *EmploymentPeriodRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM employment_periods"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете периодов работы: %w", err)
	}
	return count, nil
}

func (r *EmploymentPeriodRepositoryImpl) Create(ctx context.Context, period *models.EmploymentPeriod) (int, error) {
	query := `
		INSERT INTO employment_periods (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		period.Name,
		period.Description,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании периода работы: %w", err)
	}

	return id, nil
}

func (r *EmploymentPeriodRepositoryImpl) Update(ctx context.Context, period *models.EmploymentPeriod) error {
	query := `
		UPDATE employment_periods
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		period.Name,
		period.Description,
		period.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении периода работы: %w", err)
	}

	return nil
}

func (r *EmploymentPeriodRepositoryImpl) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := "SELECT COUNT(*) FROM reviews WHERE employment_period_id = $1"
	err := r.postgres.GetContext(ctx, &count, checkQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при проверке использования периода работы: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("период работы используется в %d отзывах и не может быть удален", count)
	}

	query := "DELETE FROM employment_periods WHERE id = $1"
	_, err = r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении периода работы: %w", err)
	}

	return nil
}
