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

type CompanyRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewCompanyRepository(postgres *db.PostgreSQL) CompanyRepository {
	return &CompanyRepositoryImpl{
		postgres: postgres,
	}
}

func (r *CompanyRepositoryImpl) Create(ctx context.Context, company *models.Company) (int, error) {
	tx, err := r.postgres.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO companies 
		(name, slug, size, logo, website, email, phone, address, city_id, reviews_count, average_rating, created_at, updated_at)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	var id int
	err = tx.QueryRowContext(
		ctx,
		query,
		company.Name,
		company.Slug,
		company.Size,
		company.Logo,
		company.Website,
		company.Email,
		company.Phone,
		company.Address,
		company.CityID,
		company.ReviewsCount,
		company.AverageRating,
		company.CreatedAt,
		company.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании компании: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}

	return id, nil
}

func (r *CompanyRepositoryImpl) GetByID(ctx context.Context, id int) (*models.CompanyWithRatings, error) {
	query := `
		SELECT id, name, slug, size, logo, website, email, phone, address, city_id,
		       reviews_count, average_rating, recommendation_percentage, created_at, updated_at
		FROM companies 
		WHERE id = $1
	`

	var company models.Company
	err := r.postgres.GetContext(ctx, &company, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("компания не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении компании: %w", err)
	}

	categoryRatings, err := r.GetCategoryRatings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рейтингов компании: %w", err)
	}

	industriesRepo := NewIndustryRepository(r.postgres)
	industries, err := industriesRepo.GetByCompanyID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отраслей компании: %w", err)
	}

	result := &models.CompanyWithRatings{
		Company:         company,
		CategoryRatings: categoryRatings,
		Industries:      industries,
	}

	if company.CityID != nil {
		cityRepo := NewCityRepository(r.postgres)
		city, err := cityRepo.GetByID(ctx, *company.CityID)
		if err == nil {
			result.City = city
		}
	}

	return result, nil
}

func (r *CompanyRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.CompanyWithRatings, error) {
	query := `
		SELECT id, name, slug, size, logo, website, email, phone, address, city_id,
		       reviews_count, average_rating, recommendation_percentage, created_at, updated_at
		FROM companies 
		WHERE slug = $1
	`

	var company models.Company
	err := r.postgres.GetContext(ctx, &company, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("компания не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении компании: %w", err)
	}

	categoryRatings, err := r.GetCategoryRatings(ctx, company.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рейтингов компании: %w", err)
	}

	industriesRepo := NewIndustryRepository(r.postgres)
	industries, err := industriesRepo.GetByCompanyID(ctx, company.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отраслей компании: %w", err)
	}

	result := &models.CompanyWithRatings{
		Company:         company,
		CategoryRatings: categoryRatings,
		Industries:      industries,
	}

	if company.CityID != nil {
		cityRepo := NewCityRepository(r.postgres)
		city, err := cityRepo.GetByID(ctx, *company.CityID)
		if err == nil {
			result.City = city
		}
	}

	return result, nil
}

func (r *CompanyRepositoryImpl) GetByName(ctx context.Context, name string) (*models.Company, error) {
	query := `
		SELECT id, name, slug, size, logo, website, email, phone, address, city_id,
		       reviews_count, average_rating, recommendation_percentage, created_at, updated_at
		FROM companies 
		WHERE name = $1
	`

	var company models.Company
	err := r.postgres.GetContext(ctx, &company, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("компания не найдена")
		}
		return nil, fmt.Errorf("ошибка при получении компании: %w", err)
	}

	return &company, nil
}

func (r *CompanyRepositoryImpl) Update(ctx context.Context, company *models.Company) error {
	query := `
		UPDATE companies
		SET name = $1, slug = $2, size = $3, logo = $4, website = $5, email = $6, phone = $7, address = $8, 
		    city_id = $9, reviews_count = $10, average_rating = $11, recommendation_percentage = $12, updated_at = $13
		WHERE id = $14
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		company.Name,
		company.Slug,
		company.Size,
		company.Logo,
		company.Website,
		company.Email,
		company.Phone,
		company.Address,
		company.CityID,
		company.ReviewsCount,
		company.AverageRating,
		company.RecommendationPercent,
		company.UpdatedAt,
		company.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении компании: %w", err)
	}

	return nil
}

func (r *CompanyRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM companies
		WHERE id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении компании: %w", err)
	}

	return nil
}

func (r *CompanyRepositoryImpl) GetAll(ctx context.Context, filter models.CompanyFilter) ([]models.CompanyWithRatings, int, error) {
	fmt.Printf("Repository GetAll: Size=%s, CityID=%v, Industries=%v\n",
		filter.Size, filter.CityID, filter.Industries)

	baseQuery := `
		FROM companies c
	`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("c.name ILIKE $%d", argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	if len(filter.Industries) > 0 {
		placeholders := make([]string, len(filter.Industries))
		for i, industryID := range filter.Industries {
			placeholders[i] = fmt.Sprintf("$%d", argID)
			args = append(args, industryID)
			argID++
		}

		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM company_industries ci WHERE ci.company_id = c.id AND ci.industry_id IN (%s))", strings.Join(placeholders, ", ")))
	}

	if filter.Size != "" {
		conditions = append(conditions, fmt.Sprintf("c.size = $%d", argID))
		args = append(args, filter.Size)
		argID++
	}

	if filter.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("c.average_rating >= $%d", argID))
		args = append(args, *filter.MinRating)
		argID++
	}

	if filter.City != "" {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM reviews r WHERE r.company_id = c.id AND r.city ILIKE $%d)", argID))
		args = append(args, "%"+filter.City+"%")
		argID++
	}

	if filter.CityID != nil {
		conditions = append(conditions, fmt.Sprintf("c.city_id = $%d", argID))
		args = append(args, *filter.CityID)
		argID++
	}

	queryConditions := baseQuery + " WHERE 1=1 "
	if len(conditions) > 0 {
		queryConditions += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(DISTINCT c.id) " + queryConditions

	sortBy := "c.name"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name":
			sortBy = "c.name"
		case "rating":
			sortBy = "c.average_rating"
		case "reviews_count":
			sortBy = "c.reviews_count"
		case "created_at":
			sortBy = "c.created_at"
		}
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
		SELECT c.id, c.name, c.slug, c.size, c.logo, c.website, c.email, c.phone, c.address, c.city_id,
		       c.reviews_count, c.average_rating, c.recommendation_percentage, c.created_at, c.updated_at
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	fmt.Printf("SQL запрос: %s\nПараметры: %v\n", dataQuery, args)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете компаний: %w", err)
	}

	var companies []models.Company
	err = r.postgres.SelectContext(ctx, &companies, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении компаний: %w", err)
	}

	result := make([]models.CompanyWithRatings, len(companies))
	industriesRepo := NewIndustryRepository(r.postgres)
	cityRepo := NewCityRepository(r.postgres)

	for i, company := range companies {
		categoryRatings, err := r.GetCategoryRatings(ctx, company.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении рейтингов компании: %w", err)
		}

		industries, err := industriesRepo.GetByCompanyID(ctx, company.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении отраслей компании: %w", err)
		}

		result[i] = models.CompanyWithRatings{
			Company:         company,
			CategoryRatings: categoryRatings,
			Industries:      industries,
		}

		if company.CityID != nil {
			city, err := cityRepo.GetByID(ctx, *company.CityID)
			if err == nil {
				result[i].City = city
			}
		}
	}

	return result, total, nil
}

func (r *CompanyRepositoryImpl) UpdateRating(ctx context.Context, companyID int) error {
	tx, err := r.postgres.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()

	updateRatingQuery := `
		UPDATE companies
		SET average_rating = COALESCE((
			SELECT AVG(rating)
			FROM reviews
			WHERE company_id = $1 AND status = 'approved'
		), 0),
		reviews_count = (
			SELECT COUNT(*)
			FROM reviews
			WHERE company_id = $1 AND status = 'approved'
		),
		recommendation_percentage = COALESCE((
			SELECT (SUM(CASE WHEN is_recommended THEN 1 ELSE 0 END) * 100.0 / COUNT(*))
			FROM reviews
			WHERE company_id = $1 AND status = 'approved'
		), 0),
		updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.Exec(updateRatingQuery, companyID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении рейтинга компании: %w", err)
	}

	deleteRatingsQuery := `
		DELETE FROM company_category_ratings
		WHERE company_id = $1
	`
	_, err = tx.Exec(deleteRatingsQuery, companyID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении рейтингов компании по категориям: %w", err)
	}

	insertRatingsQuery := `
		INSERT INTO company_category_ratings (company_id, category_id, rating)
		SELECT r.company_id, rcr.category_id, AVG(rcr.rating)
		FROM reviews r
		JOIN review_category_ratings rcr ON r.id = rcr.review_id
		WHERE r.company_id = $1 AND r.status = 'approved'
		GROUP BY r.company_id, rcr.category_id
	`
	_, err = tx.Exec(insertRatingsQuery, companyID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении рейтингов компании по категориям: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}

	return nil
}

func (r *CompanyRepositoryImpl) AddCategoryRating(ctx context.Context, companyID int, categoryID int, rating float64) error {
	query := `
		INSERT INTO company_category_ratings (company_id, category_id, rating)
		VALUES ($1, $2, $3)
		ON CONFLICT (company_id, category_id) DO UPDATE SET rating = $3
	`

	_, err := r.postgres.ExecContext(ctx, query, companyID, categoryID, rating)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении рейтинга компании по категории: %w", err)
	}

	return nil
}

func (r *CompanyRepositoryImpl) GetCategoryRatings(ctx context.Context, companyID int) ([]models.CompanyCategoryRating, error) {
	query := `
		SELECT ccr.company_id, ccr.category_id, rc.name AS category, ccr.rating
		FROM company_category_ratings ccr
		JOIN rating_categories rc ON ccr.category_id = rc.id
		WHERE ccr.company_id = $1
	`

	var ratings []models.CompanyCategoryRating
	err := r.postgres.SelectContext(ctx, &ratings, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рейтингов компании по категориям: %w", err)
	}

	return ratings, nil
}
