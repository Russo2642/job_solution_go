package handlers

import (
	"net/http"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type EmploymentPeriodHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewEmploymentPeriodHandler(postgres *db.PostgreSQL, cfg *config.Config) *EmploymentPeriodHandler {
	repo := repository.NewRepository(postgres)
	return &EmploymentPeriodHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// @Summary Получение всех периодов работы
// @Description Возвращает список всех доступных периодов работы
// @Tags employment-periods
// @Accept json
// @Produce json
// @Success 200 {object} utils.ResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /employment-periods [get]
func (h *EmploymentPeriodHandler) GetAll(c *gin.Context) {
	employmentPeriods, err := h.repo.EmploymentPeriods.GetAll(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении периодов работы", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"employment_periods": employmentPeriods,
	})
}

// @Summary Получение периода работы по ID
// @Description Возвращает информацию о периоде работы по его ID
// @Tags employment-periods
// @Accept json
// @Produce json
// @Param id path int true "ID периода работы"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /employment-periods/{id} [get]
func (h *EmploymentPeriodHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	employmentPeriod, err := h.repo.EmploymentPeriods.GetByID(c, id)
	if err != nil {
		if err.Error() == "ошибка при получении периода работы по ID: sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, "Период работы не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении периода работы", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, employmentPeriod)
}
