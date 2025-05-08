package models

import (
	"database/sql"
	"math"
	"time"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type Review struct {
	ID                 int            `json:"id" db:"id"`
	UserID             int            `json:"user_id" db:"user_id"`
	CompanyID          int            `json:"company_id" db:"company_id"`
	Position           string         `json:"position" db:"position"`
	EmploymentTypeID   *int           `json:"employment_type_id,omitempty" db:"employment_type_id"`
	EmploymentPeriodID *int           `json:"employment_period_id,omitempty" db:"employment_period_id"`
	CityID             *int           `json:"city_id,omitempty" db:"city_id"`
	Rating             float64        `json:"rating" db:"rating"`
	Pros               string         `json:"pros" db:"pros"`
	Cons               string         `json:"cons" db:"cons"`
	IsFormerEmployee   bool           `json:"is_former_employee" db:"is_former_employee"`
	Status             ReviewStatus   `json:"status" db:"status"`
	ModerationComment  sql.NullString `json:"moderation_comment,omitempty" db:"moderation_comment"`
	UsefulCount        int            `json:"useful_count" db:"useful_count"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
	ApprovedAt         sql.NullTime   `json:"approved_at,omitempty" db:"approved_at"`
}

type RatingCategory struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
}

type ReviewCategoryRating struct {
	ReviewID   int     `json:"-" db:"review_id"`
	CategoryID int     `json:"category_id" db:"category_id"`
	Category   string  `json:"category" db:"category"`
	Rating     float64 `json:"rating" db:"rating"`
}

type BenefitType struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
}

type ReviewBenefit struct {
	ID            int    `json:"-" db:"id"`
	ReviewID      int    `json:"-" db:"review_id"`
	BenefitTypeID int    `json:"benefit_type_id" db:"benefit_type_id"`
	Benefit       string `json:"benefit" db:"benefit"`
}

type UsefulMark struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ReviewID  int       `json:"review_id" db:"review_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ReviewWithDetails struct {
	Review           Review                 `json:"review"`
	CategoryRatings  []ReviewCategoryRating `json:"category_ratings"`
	Benefits         []ReviewBenefit        `json:"benefits"`
	Company          *CompanyWithRatings    `json:"company,omitempty"`
	User             *User                  `json:"user,omitempty"`
	City             *City                  `json:"city,omitempty"`
	EmploymentType   *EmploymentType        `json:"employment_type,omitempty"`
	EmploymentPeriod *EmploymentPeriod      `json:"employment_period,omitempty"`
	IsMarkedAsUseful bool                   `json:"is_marked_as_useful"`
}

type ReviewInput struct {
	CompanyID          int             `json:"company_id" binding:"required,min=1"`
	Position           string          `json:"position" binding:"required,min=2,max=100"`
	EmploymentTypeID   int             `json:"employment_type_id" binding:"required,min=1"`
	EmploymentPeriodID int             `json:"employment_period_id" binding:"required,min=1"`
	CityID             int             `json:"city_id" binding:"required,min=1"`
	CategoryRatings    map[int]float64 `json:"category_ratings" binding:"required,min=1,dive,min=1,max=5"`
	Pros               string          `json:"pros" binding:"required,min=10"`
	Cons               string          `json:"cons" binding:"required,min=10"`
	BenefitTypeIDs     []int           `json:"benefit_type_ids" binding:"omitempty,dive,min=1"`
	IsFormerEmployee   bool            `json:"is_former_employee"`
}

type ReviewFilter struct {
	CompanyID        *int          `form:"company_id" binding:"omitempty,min=1"`
	UserID           *int          `form:"user_id" binding:"omitempty,min=1"`
	Status           *ReviewStatus `form:"status" binding:"omitempty,oneof=pending approved rejected"`
	CityID           *int          `form:"city_id" binding:"omitempty,min=1"`
	MinRating        *float64      `form:"min_rating" binding:"omitempty,min=1,max=5"`
	MaxRating        *float64      `form:"max_rating" binding:"omitempty,min=1,max=5"`
	IsFormerEmployee *bool         `form:"is_former_employee" binding:"omitempty"`
	SortBy           string        `form:"sort_by" binding:"omitempty,oneof=rating created_at useful_count"`
	SortOrder        string        `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page             int           `form:"page" binding:"omitempty,min=1"`
	Limit            int           `form:"limit" binding:"omitempty,min=1,max=100"`
}

type ReviewModerationInput struct {
	Status            ReviewStatus `json:"status" binding:"required,oneof=approved rejected"`
	ModerationComment string       `json:"moderation_comment" binding:"omitempty"`
}

func NewReview(userID int, input ReviewInput) *Review {
	now := time.Now()

	var totalRating float64
	var count int
	for _, rating := range input.CategoryRatings {
		totalRating += rating
		count++
	}

	averageRating := 0.0
	if count > 0 {
		averageRating = totalRating / float64(count)
		averageRating = math.Round(averageRating*10) / 10
	}

	return &Review{
		UserID:             userID,
		CompanyID:          input.CompanyID,
		Position:           input.Position,
		EmploymentTypeID:   &input.EmploymentTypeID,
		EmploymentPeriodID: &input.EmploymentPeriodID,
		CityID:             &input.CityID,
		Rating:             averageRating,
		Pros:               input.Pros,
		Cons:               input.Cons,
		IsFormerEmployee:   input.IsFormerEmployee,
		Status:             ReviewStatusPending,
		UsefulCount:        0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (r *Review) ApproveReview(comment string) {
	now := time.Now()
	r.Status = ReviewStatusApproved
	if comment != "" {
		r.ModerationComment = sql.NullString{String: comment, Valid: true}
	}
	r.UpdatedAt = now
	r.ApprovedAt = sql.NullTime{Time: now, Valid: true}
}

func (r *Review) RejectReview(comment string) {
	r.Status = ReviewStatusRejected
	if comment != "" {
		r.ModerationComment = sql.NullString{String: comment, Valid: true}
	}
	r.UpdatedAt = time.Now()
}
