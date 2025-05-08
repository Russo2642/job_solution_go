package handlers

import (
	"net/http"
	"time"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/middleware"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewReviewHandler(postgres *db.PostgreSQL, cfg *config.Config) *ReviewHandler {
	repo := repository.NewRepository(postgres)
	return &ReviewHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// CreateReview создает новый отзыв о компании
// @Summary Создание отзыва
// @Description Создает новый отзыв о компании
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.ReviewInput true "Данные отзыва"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	_, err := h.repo.Companies.GetByID(c, input.CompanyID)
	if err != nil {
		if err.Error() == "компания не найдена" {
			utils.ErrorResponse(c, http.StatusNotFound, "Компания не найдена", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке компании", err)
		}
		return
	}

	_, err = h.repo.Cities.GetByID(c, input.CityID)
	if err != nil {
		if err.Error() == "город не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Указанный город не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке города", err)
		}
		return
	}

	_, err = h.repo.EmploymentPeriods.GetByID(c, input.EmploymentPeriodID)
	if err != nil {
		if err.Error() == "ошибка при получении периода работы по ID: sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, "Указанный период работы не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке периода работы", err)
		}
		return
	}

	_, err = h.repo.EmploymentTypes.GetByID(c, input.EmploymentTypeID)
	if err != nil {
		if err.Error() == "ошибка при получении типа занятости по ID: sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, "Указанный тип занятости не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке типа занятости", err)
		}
		return
	}

	review := models.NewReview(userID.(int), input)

	id, err := h.repo.Reviews.Create(c, review)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении отзыва", err)
		return
	}

	for categoryID, rating := range input.CategoryRatings {
		err = h.repo.Reviews.AddCategoryRating(c, id, categoryID, rating)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении рейтингов по категориям", err)
			return
		}
	}

	for _, benefitTypeID := range input.BenefitTypeIDs {
		err = h.repo.Reviews.AddBenefit(c, id, benefitTypeID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении льгот", err)
			return
		}
	}

	review.ID = id

	utils.Response(c, http.StatusCreated, gin.H{
		"review": review,
		"status": "Отзыв отправлен на модерацию",
	})
}

// GetReview возвращает отзыв по ID
// @Summary Получение отзыва
// @Description Возвращает отзыв по его ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	review, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		if err.Error() == "отзыв не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден или ожидает модерации", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзыва", err)
		}
		return
	}

	if review.Review.Status != models.ReviewStatusApproved {
		utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден или ожидает модерации", nil)
		return
	}

	utils.Response(c, http.StatusOK, review)
}

// GetCompanyReviews возвращает отзывы о компании
// @Summary Отзывы о компании
// @Description Возвращает список отзывов о компании
// @Tags reviews
// @Accept json
// @Produce json
// @Param companyId path int true "ID компании"
// @Param sort_by query string false "Поле для сортировки (rating, created_at)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Param city_id query int false "Фильтр по ID города"
// @Param min_rating query number false "Минимальный рейтинг (от 1 до 5)"
// @Param max_rating query number false "Максимальный рейтинг (от 1 до 5)"
// @Param is_former_employee query boolean false "Фильтр по статусу бывшего сотрудника (true/false)"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/company/{companyId} [get]
func (h *ReviewHandler) GetCompanyReviews(c *gin.Context) {
	companyID, err := utils.ParseIDParam(c, "companyId")
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

	var filter models.ReviewFilter
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
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	status := models.ReviewStatusApproved
	filter.Status = &status

	filter.CompanyID = &companyID

	reviews, total, err := h.repo.Reviews.GetByCompany(c, companyID, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзывов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"reviews": reviews,
		"pagination": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// GetPendingReviews возвращает отзывы, ожидающие модерации
// @Summary Отзывы на модерации
// @Description Возвращает список отзывов, ожидающих модерации
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/moderation/pending [get]
func (h *ReviewHandler) GetPendingReviews(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || (roleValue.(models.UserRole) != models.RoleModerator && roleValue.(models.UserRole) != models.RoleAdmin) {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав для просмотра отзывов на модерации", nil)
		return
	}

	var filter models.ReviewFilter
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
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	status := models.ReviewStatusPending
	filter.Status = &status

	reviews, total, err := h.repo.Reviews.GetPending(c, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзывов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"reviews": reviews,
		"pagination": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// ApproveReview одобряет отзыв
// @Summary Одобрение отзыва
// @Description Одобряет отзыв, прошедший модерацию
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Param input body models.ReviewModerationInput true "Данные модерации"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/{id}/approve [put]
func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || (roleValue.(models.UserRole) != models.RoleModerator && roleValue.(models.UserRole) != models.RoleAdmin) {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав для модерации отзывов", nil)
		return
	}

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	var input models.ReviewModerationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	if input.Status != models.ReviewStatusApproved {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный статус модерации", nil)
		return
	}

	reviewDetails, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		if err.Error() == "отзыв не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзыва", err)
		}
		return
	}

	review := reviewDetails.Review

	if review.Status != models.ReviewStatusPending {
		utils.ErrorResponse(c, http.StatusBadRequest, "Отзыв уже прошел модерацию", nil)
		return
	}

	now := time.Now()
	review.Status = models.ReviewStatusApproved
	if input.ModerationComment != "" {
		review.ModerationComment.String = input.ModerationComment
		review.ModerationComment.Valid = true
	}
	review.UpdatedAt = now
	review.ApprovedAt.Time = now
	review.ApprovedAt.Valid = true

	if err := h.repo.Reviews.Update(c, &review); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении отзыва", err)
		return
	}

	if err := h.repo.Companies.UpdateRating(c, review.CompanyID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении рейтинга компании", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Отзыв успешно одобрен",
		"review":  review,
	})
}

// RejectReview отклоняет отзыв
// @Summary Отклонение отзыва
// @Description Отклоняет отзыв с указанием причины
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Param input body models.ReviewModerationInput true "Данные модерации"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/{id}/reject [put]
func (h *ReviewHandler) RejectReview(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || (roleValue.(models.UserRole) != models.RoleModerator && roleValue.(models.UserRole) != models.RoleAdmin) {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав для модерации отзывов", nil)
		return
	}

	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	var input models.ReviewModerationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	if input.Status != models.ReviewStatusRejected {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный статус модерации", nil)
		return
	}

	if input.ModerationComment == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Необходимо указать причину отклонения отзыва", nil)
		return
	}

	reviewDetails, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		if err.Error() == "отзыв не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзыва", err)
		}
		return
	}

	review := reviewDetails.Review

	if review.Status != models.ReviewStatusPending {
		utils.ErrorResponse(c, http.StatusBadRequest, "Отзыв уже прошел модерацию", nil)
		return
	}

	review.Status = models.ReviewStatusRejected
	review.ModerationComment.String = input.ModerationComment
	review.ModerationComment.Valid = true
	review.UpdatedAt = time.Now()

	if err := h.repo.Reviews.Update(c, &review); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении отзыва", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Отзыв отклонен",
		"review":  review,
	})
}

// MarkReviewAsUseful отмечает отзыв как полезный
// @Summary Отметить отзыв как полезный
// @Description Добавляет отметку "полезно" для отзыва
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/{id}/useful [post]
func (h *ReviewHandler) MarkReviewAsUseful(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	review, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		if err.Error() == "отзыв не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзыва", err)
		}
		return
	}

	if review.Review.Status != models.ReviewStatusApproved {
		utils.ErrorResponse(c, http.StatusBadRequest, "Нельзя отметить как полезный неодобренный отзыв", nil)
		return
	}

	isMarked, err := h.repo.Reviews.HasUserMarkedReviewAsUseful(c, userID.(int), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке наличия отметки", err)
		return
	}

	if isMarked {
		utils.ErrorResponse(c, http.StatusBadRequest, "Вы уже отметили этот отзыв как полезный", nil)
		return
	}

	if err := h.repo.Reviews.AddUsefulMark(c, userID.(int), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при добавлении отметки 'полезно'", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Отзыв отмечен как полезный",
	})
}

// RemoveUsefulMark удаляет отметку "полезно" с отзыва
// @Summary Убрать отметку "полезно"
// @Description Удаляет отметку "полезно" с отзыва
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /reviews/{id}/useful [delete]
func (h *ReviewHandler) RemoveUsefulMark(c *gin.Context) {
	id, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	_, err = h.repo.Reviews.GetByID(c, id)
	if err != nil {
		if err.Error() == "отзыв не найден" {
			utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении отзыва", err)
		}
		return
	}

	isMarked, err := h.repo.Reviews.HasUserMarkedReviewAsUseful(c, userID.(int), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке наличия отметки", err)
		return
	}

	if !isMarked {
		utils.ErrorResponse(c, http.StatusBadRequest, "Вы не отмечали этот отзыв как полезный", nil)
		return
	}

	if err := h.repo.Reviews.RemoveUsefulMark(c, userID.(int), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении отметки 'полезно'", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Отметка 'полезно' удалена",
	})
}
