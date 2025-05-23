package models

import (
	"time"
)

type Company struct {
	ID                    int       `json:"id" db:"id"`
	Name                  string    `json:"name" db:"name"`
	Slug                  string    `json:"slug" db:"slug"`
	Size                  string    `json:"size" db:"size"`
	Logo                  string    `json:"logo,omitempty" db:"logo"`
	Website               string    `json:"website,omitempty" db:"website"`
	Email                 string    `json:"email,omitempty" db:"email"`
	Phone                 string    `json:"phone,omitempty" db:"phone"`
	Address               string    `json:"address,omitempty" db:"address"`
	CityID                *int      `json:"city_id,omitempty" db:"city_id"`
	ReviewsCount          int       `json:"reviews_count" db:"reviews_count"`
	AverageRating         float64   `json:"average_rating" db:"average_rating"`
	RecommendationPercent float64   `json:"recommendation_percentage" db:"recommendation_percentage"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

type CompanyCategoryRating struct {
	CompanyID  int     `json:"-" db:"company_id"`
	CategoryID int     `json:"category_id" db:"category_id"`
	Category   string  `json:"category" db:"category"`
	Rating     float64 `json:"rating" db:"rating"`
}

type CompanyWithRatings struct {
	Company         Company                 `json:"company"`
	CategoryRatings []CompanyCategoryRating `json:"category_ratings"`
	Industries      []Industry              `json:"industries"`
	City            *City                   `json:"city,omitempty"`
}

type CompanyInput struct {
	Name       string `json:"name" binding:"required,min=2,max=255"`
	Size       string `json:"size" binding:"required,oneof=small medium large enterprise"`
	Logo       string `json:"logo,omitempty" binding:"omitempty,url"`
	Website    string `json:"website,omitempty" binding:"omitempty,url"`
	Email      string `json:"email,omitempty" binding:"omitempty,email"`
	Phone      string `json:"phone,omitempty" binding:"omitempty"`
	Address    string `json:"address,omitempty" binding:"omitempty"`
	CityID     *int   `json:"city_id,omitempty" binding:"omitempty,min=1"`
	Industries []int  `json:"industries" binding:"required,min=1,dive,min=1"`
}

type CompanyFilter struct {
	Search     string   `form:"search" binding:"omitempty"`
	Industries []int    `form:"industries" binding:"omitempty,dive,min=1"`
	Size       string   `form:"size" binding:"omitempty,oneof=small medium large enterprise"`
	MinRating  *float64 `form:"min_rating" binding:"omitempty,min=1,max=5"`
	City       string   `form:"city" binding:"omitempty"`
	CityID     *int     `form:"city_id" binding:"omitempty,min=1"`
	SortBy     string   `form:"sort_by" binding:"omitempty,oneof=name rating reviews_count created_at"`
	SortOrder  string   `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page       int      `form:"page" binding:"omitempty,min=1"`
	Limit      int      `form:"limit" binding:"omitempty,min=1,max=100"`
}

var CompanySizes = map[string]string{
	"small":      "до 50 сотрудников",
	"medium":     "50-200 сотрудников",
	"large":      "200-1000 сотрудников",
	"enterprise": "более 1000 сотрудников",
}

func NewCompany(input CompanyInput) *Company {
	now := time.Now()

	tempSlug := input.Name

	return &Company{
		Name:          input.Name,
		Slug:          tempSlug,
		Size:          input.Size,
		Logo:          input.Logo,
		Website:       input.Website,
		Email:         input.Email,
		Phone:         input.Phone,
		Address:       input.Address,
		CityID:        input.CityID,
		ReviewsCount:  0,
		AverageRating: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
