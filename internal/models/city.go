package models

type City struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Region  string `json:"region" db:"region"`
	Country string `json:"country" db:"country"`
}

type CityFilter struct {
	Search    string `form:"search" binding:"omitempty"`
	Country   string `form:"country" binding:"omitempty"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name region"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
