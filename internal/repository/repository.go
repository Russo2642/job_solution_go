package repository

import (
	"context"
	"job_solition/internal/db"
	"job_solition/internal/models"
)

type Repository struct {
	Users               UserRepository
	Companies           CompanyRepository
	Reviews             ReviewRepository
	RefreshTokens       RefreshTokenRepository
	PasswordResetTokens PasswordResetRepository
	Cities              CityRepository
	Industries          IndustryRepository
	RatingCategories    RatingCategoryRepository
	BenefitTypes        BenefitTypeRepository
	EmploymentPeriods   EmploymentPeriodRepository
	EmploymentTypes     EmploymentTypeRepository
	Suggestions         SuggestionRepository
}

func NewRepository(postgres *db.PostgreSQL) *Repository {
	return &Repository{
		Users:               NewUserRepository(postgres),
		Companies:           NewCompanyRepository(postgres),
		Reviews:             NewReviewRepository(postgres),
		RefreshTokens:       NewRefreshTokenRepository(postgres),
		PasswordResetTokens: NewPasswordResetRepository(postgres),
		Cities:              NewCityRepository(postgres),
		Industries:          NewIndustryRepository(postgres),
		RatingCategories:    NewRatingCategoryRepository(postgres),
		BenefitTypes:        NewBenefitTypeRepository(postgres),
		EmploymentPeriods:   NewEmploymentPeriodRepository(postgres),
		EmploymentTypes:     NewEmploymentTypeRepository(postgres),
		Suggestions:         NewSuggestionRepository(postgres),
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	GetAll(ctx context.Context, page, limit int) ([]models.User, int, error)
	CountByRole(ctx context.Context, role models.UserRole) (int, error)
}

type CompanyRepository interface {
	Create(ctx context.Context, company *models.Company) (int, error)
	GetByID(ctx context.Context, id int) (*models.CompanyWithRatings, error)
	GetBySlug(ctx context.Context, slug string) (*models.CompanyWithRatings, error)
	GetByName(ctx context.Context, name string) (*models.Company, error)
	GetAll(ctx context.Context, filter models.CompanyFilter) ([]models.CompanyWithRatings, int, error)
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id int) error
	UpdateRating(ctx context.Context, companyID int) error
	AddCategoryRating(ctx context.Context, companyID int, categoryID int, rating float64) error
	GetCategoryRatings(ctx context.Context, companyID int) ([]models.CompanyCategoryRating, error)
	Count(ctx context.Context) (int, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) (int, error)
	GetByID(ctx context.Context, id int) (*models.ReviewWithDetails, error)
	GetByCompany(ctx context.Context, companyID int, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error)
	GetByUser(ctx context.Context, userID int, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error)
	GetPending(ctx context.Context, filter models.ReviewFilter) ([]models.ReviewWithDetails, int, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id int) error
	AddCategoryRating(ctx context.Context, reviewID int, categoryID int, rating float64) error
	AddBenefit(ctx context.Context, reviewID int, benefitTypeID int) error
	GetCategoryRatings(ctx context.Context, reviewID int) ([]models.ReviewCategoryRating, error)
	GetBenefits(ctx context.Context, reviewID int) ([]models.ReviewBenefit, error)
	MarkReviewAsUseful(ctx context.Context, reviewID int) error
	AddUsefulMark(ctx context.Context, userID, reviewID int) error
	RemoveUsefulMark(ctx context.Context, userID, reviewID int) error
	HasUserMarkedReviewAsUseful(ctx context.Context, userID, reviewID int) (bool, error)
	GetUsefulMarksByReviews(ctx context.Context, userID int, reviewIDs []int) (map[int]bool, error)
	Count(ctx context.Context) (int, error)
	CountPending(ctx context.Context) (int, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) (int, error)
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID int) error
}

type CityRepository interface {
	GetAll(ctx context.Context, filter models.CityFilter) ([]models.City, int, error)
	GetByID(ctx context.Context, id int) (*models.City, error)
	Search(ctx context.Context, query string) ([]models.City, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, city *models.City) (int, error)
	Update(ctx context.Context, city *models.City) error
	Delete(ctx context.Context, id int) error
}

type IndustryRepository interface {
	GetAll(ctx context.Context, filter models.IndustryFilter) ([]models.Industry, int, error)
	GetByID(ctx context.Context, id int) (*models.Industry, error)
	GetByIDs(ctx context.Context, ids []int) ([]models.Industry, error)
	GetByCompanyID(ctx context.Context, companyID int) ([]models.Industry, error)
	AddCompanyIndustry(ctx context.Context, companyID, industryID int) error
	RemoveCompanyIndustry(ctx context.Context, companyID, industryID int) error
	UpdateColor(ctx context.Context, id int, color string) error
	Count(ctx context.Context) (int, error)
	GetByName(ctx context.Context, name string) (*models.Industry, error)
	Create(ctx context.Context, industry *models.Industry) (int, error)
	Update(ctx context.Context, industry *models.Industry) error
	Delete(ctx context.Context, id int) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) (int, error)
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID int) error
}

type RatingCategoryRepository interface {
	GetAll(ctx context.Context) ([]models.RatingCategory, error)
	GetByID(ctx context.Context, id int) (*models.RatingCategory, error)
	GetByName(ctx context.Context, name string) (*models.RatingCategory, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, category *models.RatingCategory) (int, error)
	Update(ctx context.Context, category *models.RatingCategory) error
	Delete(ctx context.Context, id int) error
}

type BenefitTypeRepository interface {
	GetAll(ctx context.Context) ([]models.BenefitType, error)
	GetByID(ctx context.Context, id int) (*models.BenefitType, error)
	GetByName(ctx context.Context, name string) (*models.BenefitType, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, benefitType *models.BenefitType) (int, error)
	Update(ctx context.Context, benefitType *models.BenefitType) error
	Delete(ctx context.Context, id int) error
}

type EmploymentPeriodRepository interface {
	GetAll(ctx context.Context) ([]models.EmploymentPeriod, error)
	GetByID(ctx context.Context, id int) (*models.EmploymentPeriod, error)
	GetByName(ctx context.Context, name string) (*models.EmploymentPeriod, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, period *models.EmploymentPeriod) (int, error)
	Update(ctx context.Context, period *models.EmploymentPeriod) error
	Delete(ctx context.Context, id int) error
}

type EmploymentTypeRepository interface {
	GetAll(ctx context.Context) ([]models.EmploymentType, error)
	GetByID(ctx context.Context, id int) (*models.EmploymentType, error)
	GetByName(ctx context.Context, name string) (*models.EmploymentType, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, employmentType *models.EmploymentType) (int, error)
	Update(ctx context.Context, employmentType *models.EmploymentType) error
	Delete(ctx context.Context, id int) error
}

type SuggestionRepository interface {
	Create(ctx context.Context, suggestion *models.Suggestion) (int, error)
	GetAll(ctx context.Context, filter models.SuggestionFilter) ([]models.Suggestion, int, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}
