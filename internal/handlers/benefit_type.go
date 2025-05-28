package handlers

import (
	"net/http"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type BenefitTypeHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewBenefitTypeHandler(postgres *db.PostgreSQL, cfg *config.Config) *BenefitTypeHandler {
	repo := repository.NewRepository(postgres)
	return &BenefitTypeHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// @Summary Получение всех типов бенефитов
// @Description Возвращает список всех доступных типов бенефитов
// @Tags benefit-types
// @Accept json
// @Produce json
// @Success 200 {object} utils.ResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /benefit-types [get]
func (h *BenefitTypeHandler) GetAll(c *gin.Context) {
	benefitTypes, err := h.repo.BenefitTypes.GetAll(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении типов бенефитов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"benefit_types": benefitTypes,
	})
}

// @Summary Получение типа бенефита по ID
// @Description Возвращает информацию о типе бенефита по его ID
// @Tags benefit-types
// @Accept json
// @Produce json
// @Param id path int true "ID типа бенефита"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /benefit-types/{id} [get]
func (h *BenefitTypeHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	benefitType, err := h.repo.BenefitTypes.GetByID(c, id)
	if err != nil {
		if err.Error() == "ошибка при получении типа бенефита по ID: sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, "Тип бенефита не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении типа бенефита", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, benefitType)
}
