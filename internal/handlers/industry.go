package handlers

import (
	"net/http"
	"strings"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/middleware"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type IndustryHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewIndustryHandler(postgres *db.PostgreSQL, cfg *config.Config) *IndustryHandler {
	repo := repository.NewRepository(postgres)
	return &IndustryHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// GetIndustries возвращает список отраслей с фильтрацией и пагинацией
// @Summary Список отраслей
// @Description Возвращает список отраслей с возможностью фильтрации и пагинации
// @Tags industries
// @Accept json
// @Produce json
// @Param search query string false "Поисковый запрос"
// @Param sort_by query string false "Поле для сортировки (name)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /industries [get]
func (h *IndustryHandler) GetIndustries(c *gin.Context) {
	var filter models.IndustryFilter
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

	industries, total, err := h.repo.Industries.GetAll(c, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отраслей", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"industries": industries,
		"pagination": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// GetCompanyIndustries возвращает отрасли компании
// @Summary Отрасли компании
// @Description Возвращает список отраслей указанной компании
// @Tags industries
// @Accept json
// @Produce json
// @Param id path int true "ID компании"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /industries/company/{id} [get]
func (h *IndustryHandler) GetCompanyIndustries(c *gin.Context) {
	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	_, err = h.repo.Companies.GetByID(c, companyID)
	if err != nil {
		if err.Error() == "компания не найдена" {
			utils.ErrorResponse(c, http.StatusNotFound, "Компания не найдена", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке компании", err)
		}
		return
	}

	industries, err := h.repo.Industries.GetByCompanyID(c, companyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отраслей компании", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"industries": industries,
	})
}

type UpdateColorInput struct {
	Color string `json:"color" binding:"required,min=4,max=7"`
}

// UpdateIndustryColor обновляет цвет индустрии
// @Summary Обновление цвета индустрии
// @Description Обновляет цвет (hex) указанной индустрии
// @Tags industries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID индустрии"
// @Param input body UpdateColorInput true "Данные для обновления цвета"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /industries/{id}/color [put]
func (h *IndustryHandler) UpdateIndustryColor(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав для изменения цвета индустрии", nil)
		return
	}

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	_, err = h.repo.Industries.GetByID(c, id)
	if err != nil {
		if err.Error() == "отрасль не найдена" {
			utils.ErrorResponse(c, http.StatusNotFound, "Индустрия не найдена", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке индустрии", err)
		}
		return
	}

	var input UpdateColorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации данных", err)
		return
	}

	if !strings.HasPrefix(input.Color, "#") {
		utils.ErrorResponse(c, http.StatusBadRequest, "Цвет должен быть в формате HEX (например, #FF5733)", nil)
		return
	}

	err = h.repo.Industries.UpdateColor(c, id, input.Color)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении цвета индустрии", err)
		return
	}

	industry, err := h.repo.Industries.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении обновленной индустрии", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message":  "Цвет индустрии успешно обновлен",
		"industry": industry,
	})
}
