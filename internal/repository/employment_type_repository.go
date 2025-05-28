package repository

import (
	"context"
	"database/sql"
	"errors"
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
		SELECT id, name, created_at, updated_at
		FROM employment_types
		WHERE name = $1
	`

	var employmentType models.EmploymentType
	err := r.postgres.GetContext(ctx, &employmentType, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("тип занятости не найден")
		}
		return nil, fmt.Errorf("ошибка при получении типа занятости: %w", err)
	}

	return &employmentType, nil
}

func (r *EmploymentTypeRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM employment_types"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете типов занятости: %w", err)
	}
	return count, nil
}

func (r *EmploymentTypeRepositoryImpl) Create(ctx context.Context, employmentType *models.EmploymentType) (int, error) {
	query := `
		INSERT INTO employment_types (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		employmentType.Name,
		employmentType.Description,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании типа занятости: %w", err)
	}

	return id, nil
}

func (r *EmploymentTypeRepositoryImpl) Update(ctx context.Context, employmentType *models.EmploymentType) error {
	query := `
		UPDATE employment_types
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		employmentType.Name,
		employmentType.Description,
		employmentType.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении типа занятости: %w", err)
	}

	return nil
}

func (r *EmploymentTypeRepositoryImpl) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := "SELECT COUNT(*) FROM reviews WHERE employment_type_id = $1"
	err := r.postgres.GetContext(ctx, &count, checkQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при проверке использования типа занятости: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("тип занятости используется в %d отзывах и не может быть удален", count)
	}

	query := "DELETE FROM employment_types WHERE id = $1"
	_, err = r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении типа занятости: %w", err)
	}

	return nil
}
