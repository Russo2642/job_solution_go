package handlers

import (
	"net/http"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type EmploymentTypeHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewEmploymentTypeHandler(postgres *db.PostgreSQL, cfg *config.Config) *EmploymentTypeHandler {
	repo := repository.NewRepository(postgres)
	return &EmploymentTypeHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// GetAll возвращает все типы занятости
// @Summary Получение всех типов занятости
// @Description Возвращает список всех доступных типов занятости
// @Tags employment-types
// @Accept json
// @Produce json
// @Success 200 {object} utils.ResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /employment-types [get]
func (h *EmploymentTypeHandler) GetAll(c *gin.Context) {
	employmentTypes, err := h.repo.EmploymentTypes.GetAll(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении типов занятости", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"employment_types": employmentTypes,
	})
}

// GetByID возвращает тип занятости по ID
// @Summary Получение типа занятости по ID
// @Description Возвращает информацию о типе занятости по его ID
// @Tags employment-types
// @Accept json
// @Produce json
// @Param id path int true "ID типа занятости"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /employment-types/{id} [get]
func (h *EmploymentTypeHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	employmentType, err := h.repo.EmploymentTypes.GetByID(c, id)
	if err != nil {
		if err.Error() == "ошибка при получении типа занятости по ID: sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, "Тип занятости не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении типа занятости", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, employmentType)
}
