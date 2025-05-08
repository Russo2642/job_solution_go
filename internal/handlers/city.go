package handlers

import (
	"net/http"
	"strconv"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewCityHandler(postgres *db.PostgreSQL, cfg *config.Config) *CityHandler {
	repo := repository.NewRepository(postgres)
	return &CityHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// GetCities возвращает список городов с фильтрацией и пагинацией
// @Summary Список городов
// @Description Возвращает список городов с возможностью фильтрации и пагинации
// @Tags cities
// @Accept json
// @Produce json
// @Param search query string false "Поисковый запрос"
// @Param country query string false "Фильтр по стране"
// @Param sort_by query string false "Поле для сортировки (name, region)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /cities [get]
func (h *CityHandler) GetCities(c *gin.Context) {
	var filter models.CityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации параметров", err)
		return
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "name"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	cities, total, err := h.repo.Cities.GetAll(c, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении городов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"cities": cities,
		"pagination": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// SearchCities ищет города по запросу
// @Summary Поиск городов
// @Description Ищет города по названию или региону для автодополнения
// @Tags cities
// @Accept json
// @Produce json
// @Param query query string true "Поисковый запрос"
// @Param limit query int false "Максимальное количество результатов (по умолчанию 20)"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /cities/search [get]
func (h *CityHandler) SearchCities(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Не указан поисковый запрос", nil)
		return
	}

	limit := 20
	limitParam := c.Query("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
		}
	}

	filter := models.CityFilter{
		Search: query,
		Limit:  limit,
		Page:   1,
	}

	cities, _, err := h.repo.Cities.GetAll(c, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при поиске городов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"cities": cities,
	})
}
