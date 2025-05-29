package handlers

import (
	"net/http"
	"strconv"

	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type SuggestionHandlers struct {
	repo *repository.Repository
}

func NewSuggestionHandlers(repo *repository.Repository) *SuggestionHandlers {
	return &SuggestionHandlers{
		repo: repo,
	}
}

// @Summary Создать новое предложение
// @Description Создает новое предложение компании или улучшение
// @Tags suggestions
// @Accept json
// @Produce json
// @Param suggestion body models.SuggestionInput true "Информация о предложении"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /suggestions [post]
func (h *SuggestionHandlers) CreateSuggestion(c *gin.Context) {
	var input models.SuggestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверные данные", err)
		return
	}

	suggestion := models.NewSuggestion(input)
	id, err := h.repo.Suggestions.Create(c.Request.Context(), suggestion)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании предложения", err)
		return
	}

	utils.Response(c, http.StatusCreated, gin.H{
		"id": id,
	})
}

// @Summary Получить все предложения
// @Description Возвращает список предложений с фильтрацией по типу и сортировкой по дате
// @Tags suggestions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Тип предложения (company/suggestion)"
// @Param sort_order query string false "Порядок сортировки (asc/desc)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /suggestions [get]
func (h *SuggestionHandlers) GetAllSuggestions(c *gin.Context) {
	var filter models.SuggestionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверные параметры запроса", err)
		return
	}

	suggestions, total, err := h.repo.Suggestions.GetAll(c.Request.Context(), filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении предложений", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"suggestions": suggestions,
		"total":       total,
		"page":        filter.Page,
		"limit":       filter.Limit,
	})
}

// @Summary Удалить предложение
// @Description Удаляет предложение по ID
// @Tags suggestions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID предложения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /suggestions/{id} [delete]
func (h *SuggestionHandlers) DeleteSuggestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный ID предложения", err)
		return
	}

	err = h.repo.Suggestions.Delete(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении предложения", err)
		return
	}

	utils.Response(c, http.StatusOK, nil)
}
