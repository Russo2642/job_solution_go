package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type BenefitTypeRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewBenefitTypeRepository(postgres *db.PostgreSQL) BenefitTypeRepository {
	return &BenefitTypeRepositoryImpl{
		postgres: postgres,
	}
}

func (r *BenefitTypeRepositoryImpl) GetAll(ctx context.Context) ([]models.BenefitType, error) {
	query := `
		SELECT id, name, description
		FROM benefit_types
		ORDER BY name
	`

	var benefitTypes []models.BenefitType
	err := r.postgres.SelectContext(ctx, &benefitTypes, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типов бенефитов: %w", err)
	}

	return benefitTypes, nil
}

func (r *BenefitTypeRepositoryImpl) GetByID(ctx context.Context, id int) (*models.BenefitType, error) {
	query := `
		SELECT id, name, description
		FROM benefit_types
		WHERE id = $1
	`

	var benefitType models.BenefitType
	err := r.postgres.GetContext(ctx, &benefitType, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типа бенефита по ID: %w", err)
	}

	return &benefitType, nil
}

func (r *BenefitTypeRepositoryImpl) GetByName(ctx context.Context, name string) (*models.BenefitType, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM benefit_types
		WHERE name = $1
	`

	var benefitType models.BenefitType
	err := r.postgres.GetContext(ctx, &benefitType, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("тип бенефита не найден")
		}
		return nil, fmt.Errorf("ошибка при получении типа бенефита: %w", err)
	}

	return &benefitType, nil
}

func (r *BenefitTypeRepositoryImpl) Count(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM benefit_types"
	err := r.postgres.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете типов бенефитов: %w", err)
	}
	return count, nil
}

func (r *BenefitTypeRepositoryImpl) Create(ctx context.Context, benefitType *models.BenefitType) (int, error) {
	query := `
		INSERT INTO benefit_types (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRowContext(
		ctx,
		query,
		benefitType.Name,
		benefitType.Description,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании типа бенефита: %w", err)
	}

	return id, nil
}

func (r *BenefitTypeRepositoryImpl) Update(ctx context.Context, benefitType *models.BenefitType) error {
	query := `
		UPDATE benefit_types
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		benefitType.Name,
		benefitType.Description,
		benefitType.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении типа бенефита: %w", err)
	}

	return nil
}

func (r *BenefitTypeRepositoryImpl) Delete(ctx context.Context, id int) error {
	var count int
	checkQuery := "SELECT COUNT(*) FROM review_benefits WHERE benefit_type_id = $1"
	err := r.postgres.GetContext(ctx, &count, checkQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при проверке использования типа бенефита: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("тип бенефита используется в отзывах и не может быть удален")
	}

	query := "DELETE FROM benefit_types WHERE id = $1"
	_, err = r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении типа бенефита: %w", err)
	}

	return nil
}
