package models

import (
	"time"
)

type SuggestionType string

const (
	SuggestionTypeCompany     SuggestionType = "company"
	SuggestionTypeImprovement SuggestionType = "suggestion"
)

type Suggestion struct {
	ID        int            `json:"id" db:"id"`
	Type      SuggestionType `json:"type" db:"type"`
	Text      string         `json:"text" db:"text"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type SuggestionInput struct {
	Type SuggestionType `json:"type" binding:"required,oneof=company suggestion"`
	Text string         `json:"text" binding:"required,min=5,max=2000"`
}

type SuggestionFilter struct {
	Type      SuggestionType `form:"type" binding:"omitempty,oneof=company suggestion"`
	SortOrder string         `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page      int            `form:"page" binding:"omitempty,min=1"`
	Limit     int            `form:"limit" binding:"omitempty,min=1,max=100"`
}

func NewSuggestion(input SuggestionInput) *Suggestion {
	return &Suggestion{
		Type:      input.Type,
		Text:      input.Text,
		CreatedAt: time.Now(),
	}
}
