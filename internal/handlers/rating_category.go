package handlers

import (
	"net/http"

	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type RatingCategoryHandler struct {
	repo *repository.Repository
}

func NewRatingCategoryHandler(repo *repository.Repository) *RatingCategoryHandler {
	return &RatingCategoryHandler{
		repo: repo,
	}
}

// @Summary Получение всех категорий рейтингов
// @Description Возвращает список всех доступных категорий рейтингов
// @Tags rating-categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.ResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /rating-categories [get]
func (h *RatingCategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.repo.RatingCategories.GetAll(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении категорий рейтингов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"categories": categories,
	})
}

// @Summary Получение категории рейтинга по ID
// @Description Возвращает категорию рейтинга по её ID
// @Tags rating-categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории рейтинга"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /rating-categories/{id} [get]
func (h *RatingCategoryHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	category, err := h.repo.RatingCategories.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Категория рейтинга не найдена", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"category": category,
	})
}
