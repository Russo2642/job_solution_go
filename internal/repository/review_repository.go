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

type ReviewRepositoryImpl struct {
	postgres *db.PostgreSQL
}

func NewReviewRepository(postgres *db.PostgreSQL) ReviewRepository {
	return &ReviewRepositoryImpl{
		postgres: postgres,
	}
}

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review *models.Review) (int, error) {
	tx, err := r.postgres.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO reviews 
		(user_id, company_id, position, employment_type_id, employment_period_id, city_id, rating, pros, cons, is_former_employee, is_recommended, status, created_at, updated_at)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	var id int
	err = tx.QueryRowx(
		query,
		review.UserID,
		review.CompanyID,
		review.Position,
		review.EmploymentTypeID,
		review.EmploymentPeriodID,
		review.CityID,
		review.Rating,
		review.Pros,
		review.Cons,
		review.IsFormerEmployee,
		review.IsRecommended,
		review.Status,
		review.CreatedAt,
		review.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании отзыва: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}

	return id, nil
}

func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id int) (*models.ReviewWithDetails, error) {
	query := `
		SELECT id, user_id, company_id, position, employment_type_id, employment_period_id,
		       city_id, rating, pros, cons, is_former_employee, is_recommended, status, moderation_comment, 
		       useful_count, created_at, updated_at, approved_at
		FROM reviews
		WHERE id = $1
	`

	var review models.Review
	err := r.postgres.GetContext(ctx, &review, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("отзыв не найден")
		}
		return nil, fmt.Errorf("ошибка при получении отзыва: %w", err)
	}

	categoryRatings, err := r.GetCategoryRatings(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рейтингов по категориям: %w", err)
	}

	benefits, err := r.GetBenefits(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении льгот: %w", err)
	}

	var city *models.City
	if review.CityID != nil {
		cityRepo := NewCityRepository(r.postgres)
		cityData, err := cityRepo.GetByID(ctx, *review.CityID)
		if err == nil {
			city = cityData
		}
	}

	var employmentType *models.EmploymentType
	if review.EmploymentTypeID != nil {
		employmentTypeRepo := NewEmploymentTypeRepository(r.postgres)
		employmentTypeData, err := employmentTypeRepo.GetByID(ctx, *review.EmploymentTypeID)
		if err == nil {
			employmentType = employmentTypeData
		}
	}

	var employmentPeriod *models.EmploymentPeriod
	if review.EmploymentPeriodID != nil {
		employmentPeriodRepo := NewEmploymentPeriodRepository(r.postgres)
		employmentPeriodData, err := employmentPeriodRepo.GetByID(ctx, *review.EmploymentPeriodID)
		if err == nil {
			employmentPeriod = employmentPeriodData
		}
	}

	result := &models.ReviewWithDetails{
		Review:           review,
		CategoryRatings:  categoryRatings,
		Benefits:         benefits,
		City:             city,
		EmploymentType:   employmentType,
		EmploymentPeriod: employmentPeriod,
		IsMarkedAsUseful: false,
	}

	userID, exists := ctx.Value("user_id").(int)
	if exists && userID > 0 {
		isMarked, err := r.HasUserMarkedReviewAsUseful(ctx, userID, id)
		if err == nil {
			result.IsMarkedAsUseful = isMarked
		}
	}

	return result, nil
}

func (r *ReviewRepositoryImpl) GetByCompany(ctx context.Context, companyID int, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error) {
	filterCopy := filter
	filterCopy.CompanyID = &companyID

	reviews, total, err := r.getReviews(ctx, filterCopy)
	if err != nil {
		return nil, 0, err
	}

	userID, exists := ctx.Value("user_id").(int)
	if exists && userID > 0 && len(reviews) > 0 {
		reviewIDs := make([]int, len(reviews))
		for i, review := range reviews {
			reviewIDs[i] = review.Review.ID
		}

		userMarks, err := r.GetUsefulMarksByReviews(ctx, userID, reviewIDs)
		if err == nil {
			for i := range reviews {
				reviews[i].IsMarkedAsUseful = userMarks[reviews[i].Review.ID]
			}
		}
	}

	return reviews, total, nil
}

func (r *ReviewRepositoryImpl) getReviews(ctx context.Context, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error) {
	baseQuery := `
		FROM reviews 
		WHERE 1=1
	`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, *filter.Status)
		argID++
	}

	if filter.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argID))
		args = append(args, *filter.CompanyID)
		argID++
	}

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, *filter.UserID)
		argID++
	}

	if filter.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argID))
		args = append(args, *filter.MinRating)
		argID++
	}

	if filter.MaxRating != nil {
		conditions = append(conditions, fmt.Sprintf("rating <= $%d", argID))
		args = append(args, *filter.MaxRating)
		argID++
	}

	if filter.IsFormerEmployee != nil {
		conditions = append(conditions, fmt.Sprintf("is_former_employee = $%d", argID))
		args = append(args, *filter.IsFormerEmployee)
		argID++
	}

	if filter.CityID != nil {
		conditions = append(conditions, fmt.Sprintf("city_id = $%d", argID))
		args = append(args, *filter.CityID)
		argID++
	}

	queryConditions := baseQuery
	if len(conditions) > 0 {
		queryConditions += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + queryConditions

	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	dataQuery := fmt.Sprintf(`
		SELECT id, user_id, company_id, position, employment_type_id, employment_period_id, city_id, rating,
		       pros, cons, is_former_employee, is_recommended, status, moderation_comment, useful_count, created_at, updated_at, approved_at
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете отзывов: %w", err)
	}

	var reviews []models.Review
	err = r.postgres.SelectContext(ctx, &reviews, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении отзывов: %w", err)
	}

	result := make([]models.ReviewWithDetails, len(reviews))
	cityRepo := NewCityRepository(r.postgres)
	employmentTypeRepo := NewEmploymentTypeRepository(r.postgres)
	employmentPeriodRepo := NewEmploymentPeriodRepository(r.postgres)
	companyRepo := NewCompanyRepository(r.postgres)

	for i, review := range reviews {
		categoryRatings, err := r.GetCategoryRatings(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении рейтингов отзыва: %w", err)
		}

		benefits, err := r.GetBenefits(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении льгот отзыва: %w", err)
		}

		var city *models.City
		if review.CityID != nil && *review.CityID > 0 {
			cityData, err := cityRepo.GetByID(ctx, *review.CityID)
			if err == nil {
				city = cityData
			}
		}

		var employmentType *models.EmploymentType
		if review.EmploymentTypeID != nil {
			employmentTypeData, err := employmentTypeRepo.GetByID(ctx, *review.EmploymentTypeID)
			if err == nil {
				employmentType = employmentTypeData
			}
		}

		var employmentPeriod *models.EmploymentPeriod
		if review.EmploymentPeriodID != nil {
			employmentPeriodData, err := employmentPeriodRepo.GetByID(ctx, *review.EmploymentPeriodID)
			if err == nil {
				employmentPeriod = employmentPeriodData
			}
		}

		var company *models.CompanyWithRatings
		companyData, err := companyRepo.GetByID(ctx, review.CompanyID)
		if err == nil {
			company = companyData
		}

		result[i] = models.ReviewWithDetails{
			Review:           review,
			CategoryRatings:  categoryRatings,
			Benefits:         benefits,
			City:             city,
			EmploymentType:   employmentType,
			EmploymentPeriod: employmentPeriod,
			Company:          company,
			IsMarkedAsUseful: false,
		}
	}

	return result, total, nil
}

func (r *ReviewRepositoryImpl) GetByUser(ctx context.Context, userID int, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error) {
	baseQuery := `
		FROM reviews 
		WHERE user_id = $1
	`

	conditions := []string{}
	args := []interface{}{userID}
	argID := 2

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, *filter.Status)
		argID++
	}

	if filter.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argID))
		args = append(args, *filter.CompanyID)
		argID++
	}

	if filter.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argID))
		args = append(args, *filter.MinRating)
		argID++
	}

	if filter.CityID != nil && *filter.CityID > 0 {
		conditions = append(conditions, fmt.Sprintf("city_id = $%d", argID))
		args = append(args, *filter.CityID)
		argID++
	}

	queryConditions := baseQuery
	if len(conditions) > 0 {
		queryConditions += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + queryConditions

	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	dataQuery := fmt.Sprintf(`
		SELECT id, user_id, company_id, position, employment_type_id, employment_period_id, city_id, rating,
		       pros, cons, is_former_employee, is_recommended, status, moderation_comment, useful_count, created_at, updated_at, approved_at
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете отзывов: %w", err)
	}

	var reviews []models.Review
	err = r.postgres.SelectContext(ctx, &reviews, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении отзывов: %w", err)
	}

	result := make([]models.ReviewWithDetails, len(reviews))
	cityRepo := NewCityRepository(r.postgres)
	employmentTypeRepo := NewEmploymentTypeRepository(r.postgres)
	employmentPeriodRepo := NewEmploymentPeriodRepository(r.postgres)
	companyRepo := NewCompanyRepository(r.postgres)

	for i, review := range reviews {
		categoryRatings, err := r.GetCategoryRatings(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении рейтингов отзыва: %w", err)
		}

		benefits, err := r.GetBenefits(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении льгот отзыва: %w", err)
		}

		var city *models.City
		if review.CityID != nil && *review.CityID > 0 {
			cityData, err := cityRepo.GetByID(ctx, *review.CityID)
			if err == nil {
				city = cityData
			}
		}

		var employmentType *models.EmploymentType
		if review.EmploymentTypeID != nil {
			employmentTypeData, err := employmentTypeRepo.GetByID(ctx, *review.EmploymentTypeID)
			if err == nil {
				employmentType = employmentTypeData
			}
		}

		var employmentPeriod *models.EmploymentPeriod
		if review.EmploymentPeriodID != nil {
			employmentPeriodData, err := employmentPeriodRepo.GetByID(ctx, *review.EmploymentPeriodID)
			if err == nil {
				employmentPeriod = employmentPeriodData
			}
		}

		var company *models.CompanyWithRatings
		companyData, err := companyRepo.GetByID(ctx, review.CompanyID)
		if err == nil {
			company = companyData
		}

		result[i] = models.ReviewWithDetails{
			Review:           review,
			CategoryRatings:  categoryRatings,
			Benefits:         benefits,
			City:             city,
			EmploymentType:   employmentType,
			EmploymentPeriod: employmentPeriod,
			Company:          company,
			IsMarkedAsUseful: false,
		}
	}

	return result, total, nil
}

func (r *ReviewRepositoryImpl) GetPending(ctx context.Context, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error) {
	baseQuery := `
		FROM reviews 
		WHERE status = 'pending'
	`

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if filter.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argID))
		args = append(args, *filter.CompanyID)
		argID++
	}

	if filter.CityID != nil && *filter.CityID > 0 {
		conditions = append(conditions, fmt.Sprintf("city_id = $%d", argID))
		args = append(args, *filter.CityID)
		argID++
	}

	queryConditions := baseQuery
	if len(conditions) > 0 {
		queryConditions += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + queryConditions

	sortBy := "created_at"
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
		SELECT id, user_id, company_id, position, employment_type_id, employment_period_id, city_id, rating,
		       pros, cons, is_former_employee, is_recommended, status, moderation_comment, useful_count, created_at, updated_at, approved_at
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, queryConditions, sortBy, sortOrder, argID, argID+1)

	args = append(args, filter.Limit, offset)

	var total int
	err := r.postgres.GetContext(ctx, &total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при подсчете отзывов: %w", err)
	}

	var reviews []models.Review
	err = r.postgres.SelectContext(ctx, &reviews, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при получении отзывов: %w", err)
	}

	result := make([]models.ReviewWithDetails, len(reviews))
	cityRepo := NewCityRepository(r.postgres)
	employmentTypeRepo := NewEmploymentTypeRepository(r.postgres)
	employmentPeriodRepo := NewEmploymentPeriodRepository(r.postgres)
	companyRepo := NewCompanyRepository(r.postgres)

	for i, review := range reviews {
		categoryRatings, err := r.GetCategoryRatings(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении рейтингов отзыва: %w", err)
		}

		benefits, err := r.GetBenefits(ctx, review.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при получении льгот отзыва: %w", err)
		}

		var city *models.City
		if review.CityID != nil && *review.CityID > 0 {
			cityData, err := cityRepo.GetByID(ctx, *review.CityID)
			if err == nil {
				city = cityData
			}
		}

		var employmentType *models.EmploymentType
		if review.EmploymentTypeID != nil {
			employmentTypeData, err := employmentTypeRepo.GetByID(ctx, *review.EmploymentTypeID)
			if err == nil {
				employmentType = employmentTypeData
			}
		}

		var employmentPeriod *models.EmploymentPeriod
		if review.EmploymentPeriodID != nil {
			employmentPeriodData, err := employmentPeriodRepo.GetByID(ctx, *review.EmploymentPeriodID)
			if err == nil {
				employmentPeriod = employmentPeriodData
			}
		}

		var company *models.CompanyWithRatings
		companyData, err := companyRepo.GetByID(ctx, review.CompanyID)
		if err == nil {
			company = companyData
		}

		result[i] = models.ReviewWithDetails{
			Review:           review,
			CategoryRatings:  categoryRatings,
			Benefits:         benefits,
			City:             city,
			EmploymentType:   employmentType,
			EmploymentPeriod: employmentPeriod,
			Company:          company,
			IsMarkedAsUseful: false,
		}
	}

	return result, total, nil
}

func (r *ReviewRepositoryImpl) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews
		SET position = $1, employment_type_id = $2, employment_period_id = $3, city_id = $4, rating = $5,
		    pros = $6, cons = $7, is_former_employee = $8, is_recommended = $9, status = $10, moderation_comment = $11, updated_at = $12, approved_at = $13
		WHERE id = $14
	`

	_, err := r.postgres.ExecContext(
		ctx,
		query,
		review.Position,
		review.EmploymentTypeID,
		review.EmploymentPeriodID,
		review.CityID,
		review.Rating,
		review.Pros,
		review.Cons,
		review.IsFormerEmployee,
		review.IsRecommended,
		review.Status,
		review.ModerationComment,
		review.UpdatedAt,
		review.ApprovedAt,
		review.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении отзыва: %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM reviews
		WHERE id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении отзыва: %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) AddCategoryRating(ctx context.Context, reviewID int, categoryID int, rating float64) error {
	query := `
		INSERT INTO review_category_ratings (review_id, category_id, rating)
		VALUES ($1, $2, $3)
		ON CONFLICT (review_id, category_id) DO UPDATE SET rating = $3
	`

	_, err := r.postgres.ExecContext(ctx, query, reviewID, categoryID, rating)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении рейтинга отзыва по категории: %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) AddBenefit(ctx context.Context, reviewID int, benefitTypeID int) error {
	query := `
		INSERT INTO review_benefits (review_id, benefit_type_id)
		VALUES ($1, $2)
	`

	_, err := r.postgres.ExecContext(ctx, query, reviewID, benefitTypeID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении льготы к отзыву: %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) GetCategoryRatings(ctx context.Context, reviewID int) ([]models.ReviewCategoryRating, error) {
	query := `
		SELECT rcr.review_id, rcr.category_id, rc.name AS category, rcr.rating
		FROM review_category_ratings rcr
		JOIN rating_categories rc ON rcr.category_id = rc.id
		WHERE rcr.review_id = $1
	`

	var ratings []models.ReviewCategoryRating
	err := r.postgres.SelectContext(ctx, &ratings, query, reviewID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рейтингов отзыва по категориям: %w", err)
	}

	return ratings, nil
}

func (r *ReviewRepositoryImpl) GetBenefits(ctx context.Context, reviewID int) ([]models.ReviewBenefit, error) {
	query := `
		SELECT rb.id, rb.review_id, rb.benefit_type_id, bt.name AS benefit
		FROM review_benefits rb
		JOIN benefit_types bt ON rb.benefit_type_id = bt.id
		WHERE rb.review_id = $1
	`

	var benefits []models.ReviewBenefit
	err := r.postgres.SelectContext(ctx, &benefits, query, reviewID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении льгот отзыва: %w", err)
	}

	return benefits, nil
}

func (r *ReviewRepositoryImpl) MarkReviewAsUseful(ctx context.Context, reviewID int) error {
	query := `
		UPDATE reviews
		SET useful_count = useful_count + 1
		WHERE id = $1
	`

	_, err := r.postgres.ExecContext(ctx, query, reviewID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении счетчика полезности отзыва: %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) AddUsefulMark(ctx context.Context, userID, reviewID int) error {
	query := `
		INSERT INTO useful_marks (user_id, review_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, review_id) DO NOTHING
	`

	_, err := r.postgres.ExecContext(ctx, query, userID, reviewID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении отметки 'полезно': %w", err)
	}

	return nil
}

func (r *ReviewRepositoryImpl) RemoveUsefulMark(ctx context.Context, userID, reviewID int) error {
	query := `
		DELETE FROM useful_marks
		WHERE user_id = $1 AND review_id = $2
	`

	result, err := r.postgres.ExecContext(ctx, query, userID, reviewID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении отметки 'полезно': %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("отметка 'полезно' не найдена")
	}

	return nil
}

func (r *ReviewRepositoryImpl) HasUserMarkedReviewAsUseful(ctx context.Context, userID, reviewID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM useful_marks
			WHERE user_id = $1 AND review_id = $2
		)
	`

	var exists bool
	err := r.postgres.GetContext(ctx, &exists, query, userID, reviewID)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке наличия отметки 'полезно': %w", err)
	}

	return exists, nil
}

func (r *ReviewRepositoryImpl) GetUsefulMarksByReviews(ctx context.Context, userID int, reviewIDs []int) (map[int]bool, error) {
	if len(reviewIDs) == 0 {
		return make(map[int]bool), nil
	}

	placeholders := make([]string, len(reviewIDs))
	args := make([]interface{}, len(reviewIDs)+1)

	args[0] = userID
	for i, reviewID := range reviewIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = reviewID
	}

	query := fmt.Sprintf(`
		SELECT review_id
		FROM useful_marks
		WHERE user_id = $1 AND review_id IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := r.postgres.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отметок 'полезно': %w", err)
	}
	defer rows.Close()

	result := make(map[int]bool)
	for rows.Next() {
		var reviewID int
		if err := rows.Scan(&reviewID); err != nil {
			return nil, fmt.Errorf("ошибка при чтении ID отзыва: %w", err)
		}
		result[reviewID] = true
	}

	return result, nil
}
