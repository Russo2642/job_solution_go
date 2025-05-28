package models

type AdminStatistics struct {
	UsersCount     int `json:"users_count"`
	CompaniesCount int `json:"companies_count"`
	ReviewsCount   int `json:"reviews_count"`
	PendingReviews int `json:"pending_reviews"`

	CitiesCount            int `json:"cities_count"`
	IndustriesCount        int `json:"industries_count"`
	BenefitTypesCount      int `json:"benefit_types_count"`
	RatingCategoriesCount  int `json:"rating_categories_count"`
	EmploymentTypesCount   int `json:"employment_types_count"`
	EmploymentPeriodsCount int `json:"employment_periods_count"`
}

type RatingCategoryInput struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type BenefitTypeInput struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type EmploymentTypeInput struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type EmploymentPeriodInput struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type CityInput struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Region  string `json:"region,omitempty" binding:"omitempty,max=100"`
	Country string `json:"country" binding:"required,min=2,max=100"`
}

type IndustryInput struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Color string `json:"color,omitempty" binding:"omitempty,max=20"`
}

type AdminReviewUpdateInput struct {
	Position          *string  `json:"position,omitempty" binding:"omitempty,min=2,max=100"`
	Rating            *float64 `json:"rating,omitempty" binding:"omitempty,min=1,max=5"`
	Pros              *string  `json:"pros,omitempty" binding:"omitempty"`
	Cons              *string  `json:"cons,omitempty" binding:"omitempty"`
	IsFormerEmployee  *bool    `json:"is_former_employee,omitempty"`
	IsRecommended     *bool    `json:"is_recommended,omitempty"`
	Status            *string  `json:"status,omitempty" binding:"omitempty,oneof=pending approved rejected"`
	ModerationComment *string  `json:"moderation_comment,omitempty" binding:"omitempty"`
}
