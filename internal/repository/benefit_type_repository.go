package repository

import (
	"context"
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
		SELECT id, name, description
		FROM benefit_types
		WHERE name = $1
	`

	var benefitType models.BenefitType
	err := r.postgres.GetContext(ctx, &benefitType, query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении типа бенефита по имени: %w", err)
	}

	return &benefitType, nil
}
