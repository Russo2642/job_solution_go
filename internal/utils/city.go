package utils

import (
	"context"
	"fmt"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"strings"
)

func FormatCity(city *models.City) string {
	if city == nil {
		return ""
	}

	parts := []string{city.Name}

	if city.Region != "" && city.Region != city.Name {
		parts = append(parts, city.Region)
	}

	if city.Country != "" {
		parts = append(parts, city.Country)
	}

	return strings.Join(parts, ", ")
}

func GetCityName(repo *repository.Repository, cityID int) string {
	if cityID <= 0 {
		return ""
	}

	city, err := repo.Cities.GetByID(context.TODO(), cityID)
	if err != nil {
		return ""
	}

	return city.Name
}

func FormatLocationInfo(review *models.ReviewWithDetails) string {
	if review == nil {
		return ""
	}

	if review.City != nil {
		return FormatCity(review.City)
	}

	return ""
}

func GetDisplayCityName(city *models.City) string {
	if city == nil {
		return ""
	}

	if city.Country == "Казахстан" {
		return fmt.Sprintf("%s, %s", city.Name, city.Country)
	}

	return FormatCity(city)
}
