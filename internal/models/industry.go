package models

type Industry struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Color string `json:"color,omitempty" db:"color"`
}

type IndustryFilter struct {
	Search    string `form:"search" binding:"omitempty"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
