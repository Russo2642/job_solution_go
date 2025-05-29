package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"job_solition/internal/config"
	"job_solition/internal/middleware"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAdminHandler(repo *repository.Repository, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// @Summary Получение статистики
// @Description Возвращает общую статистику по сайту: количество пользователей, компаний, отзывов и т.д.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.ResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/statistics [get]
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	usersCount, err := h.repo.Users.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества пользователей", err)
		return
	}

	companiesCount, err := h.repo.Companies.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества компаний", err)
		return
	}

	reviewsCount, err := h.repo.Reviews.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества отзывов", err)
		return
	}

	pendingReviews, err := h.repo.Reviews.CountPending(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества ожидающих отзывов", err)
		return
	}

	approvedReviews, err := h.repo.Reviews.CountApproved(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества одобренных отзывов", err)
		return
	}

	rejectedReviews, err := h.repo.Reviews.CountRejected(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества отклоненных отзывов", err)
		return
	}

	citiesCount, err := h.repo.Cities.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества городов", err)
		return
	}

	industriesCount, err := h.repo.Industries.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества индустрий", err)
		return
	}

	benefitTypesCount, err := h.repo.BenefitTypes.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества типов бенефитов", err)
		return
	}

	ratingCategoriesCount, err := h.repo.RatingCategories.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества категорий рейтингов", err)
		return
	}

	employmentTypesCount, err := h.repo.EmploymentTypes.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества типов занятости", err)
		return
	}

	employmentPeriodsCount, err := h.repo.EmploymentPeriods.Count(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении количества периодов работы", err)
		return
	}

	statistics := models.AdminStatistics{
		UsersCount:             usersCount,
		CompaniesCount:         companiesCount,
		ReviewsCount:           reviewsCount,
		PendingReviews:         pendingReviews,
		ApprovedReviews:        approvedReviews,
		RejectedReviews:        rejectedReviews,
		CitiesCount:            citiesCount,
		IndustriesCount:        industriesCount,
		BenefitTypesCount:      benefitTypesCount,
		RatingCategoriesCount:  ratingCategoriesCount,
		EmploymentTypesCount:   employmentTypesCount,
		EmploymentPeriodsCount: employmentPeriodsCount,
	}

	utils.Response(c, http.StatusOK, statistics)
}

// @Summary Получение списка пользователей
// @Description Возвращает список пользователей с пагинацией
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var page, limit int

	if pageStr := c.Query("page"); pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err == nil && limitVal > 0 {
			limit = limitVal
		}
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	users, total, err := h.repo.Users.GetAll(c, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении пользователей", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// @Summary Получение пользователя
// @Description Возвращает информацию о пользователе по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	user, err := h.repo.Users.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Пользователь не найден", err)
		return
	}

	utils.Response(c, http.StatusOK, user)
}

// @Summary Обновление роли пользователя
// @Description Обновляет роль пользователя (admin, moderator, user)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Param input body models.UserRoleUpdateInput true "Новая роль пользователя"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.UserRoleUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	user, err := h.repo.Users.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Пользователь не найден", err)
		return
	}

	validRoles := map[models.UserRole]bool{
		models.RoleAdmin:     true,
		models.RoleModerator: true,
		models.RoleUser:      true,
	}

	if !validRoles[input.Role] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Некорректная роль. Допустимые значения: admin, moderator, user", nil)
		return
	}

	if user.Role == models.RoleAdmin && input.Role != models.RoleAdmin {
		adminsCount, err := h.repo.Users.CountByRole(c, models.RoleAdmin)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке количества администраторов", err)
			return
		}

		if adminsCount <= 1 {
			utils.ErrorResponse(c, http.StatusBadRequest, "Невозможно понизить последнего администратора", nil)
			return
		}
	}

	user.Role = input.Role
	user.UpdatedAt = time.Now()

	if err := h.repo.Users.Update(c, user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении роли пользователя", err)
		return
	}

	utils.Response(c, http.StatusOK, user)
}

// @Summary Удаление пользователя
// @Description Удаляет пользователя по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	currentUserID, _ := c.Get(middleware.UserIDKey)
	if currentUserID.(int) == id {
		utils.ErrorResponse(c, http.StatusBadRequest, "Невозможно удалить собственную учетную запись", nil)
		return
	}

	user, err := h.repo.Users.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Пользователь не найден", err)
		return
	}

	if user.Role == models.RoleAdmin {
		adminsCount, err := h.repo.Users.CountByRole(c, models.RoleAdmin)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке количества администраторов", err)
			return
		}

		if adminsCount <= 1 {
			utils.ErrorResponse(c, http.StatusBadRequest, "Невозможно удалить последнего администратора", nil)
			return
		}
	}

	if err := h.repo.Users.Delete(c, id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении пользователя", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}

// @Summary Создание категории рейтинга
// @Description Создает новую категорию рейтинга
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.RatingCategoryInput true "Данные категории"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/rating-categories [post]
func (h *AdminHandler) CreateRatingCategory(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.RatingCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	_, err := h.repo.RatingCategories.GetByName(c, input.Name)
	if err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Категория с таким названием уже существует", nil)
		return
	}

	category := &models.RatingCategory{
		Name:        input.Name,
		Description: input.Description,
	}

	id, err := h.repo.RatingCategories.Create(c, category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании категории", err)
		return
	}

	category.ID = id

	utils.Response(c, http.StatusCreated, category)
}

// @Summary Обновление категории рейтинга
// @Description Обновляет существующую категорию рейтинга
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Param input body models.RatingCategoryInput true "Данные категории"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/rating-categories/{id} [put]
func (h *AdminHandler) UpdateRatingCategory(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.RatingCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	category, err := h.repo.RatingCategories.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Категория не найдена", err)
		return
	}

	if category.Name != input.Name {
		existingCategory, err := h.repo.RatingCategories.GetByName(c, input.Name)
		if err == nil && existingCategory != nil && existingCategory.ID != id {
			utils.ErrorResponse(c, http.StatusBadRequest, "Категория с таким названием уже существует", nil)
			return
		}
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := h.repo.RatingCategories.Update(c, category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении категории", err)
		return
	}

	utils.Response(c, http.StatusOK, category)
}

// @Summary Удаление категории рейтинга
// @Description Удаляет категорию рейтинга по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/rating-categories/{id} [delete]
func (h *AdminHandler) DeleteRatingCategory(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.RatingCategories.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Категория не найдена", err)
		return
	}

	if err := h.repo.RatingCategories.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "категория используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении категории", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Категория успешно удалена"})
}

// @Summary Обновление отзыва
// @Description Обновляет отзыв по ID (для администратора)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Param input body models.AdminReviewUpdateInput true "Данные для обновления"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/reviews/{id} [put]
func (h *AdminHandler) UpdateReview(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.AdminReviewUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	reviewDetails, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", err)
		return
	}

	review := reviewDetails.Review

	if input.Position != nil {
		review.Position = *input.Position
	}

	if input.Rating != nil {
		review.Rating = *input.Rating
	}

	if input.Pros != nil {
		review.Pros = *input.Pros
	}

	if input.Cons != nil {
		review.Cons = *input.Cons
	}

	if input.IsFormerEmployee != nil {
		review.IsFormerEmployee = *input.IsFormerEmployee
	}

	if input.IsRecommended != nil {
		review.IsRecommended = *input.IsRecommended
	}

	if input.Status != nil {
		review.Status = models.ReviewStatus(*input.Status)

		if *input.Status == "approved" {
			now := time.Now()
			review.ApprovedAt = sql.NullTime{Time: now, Valid: true}
		}
	}

	if input.ModerationComment != nil {
		review.ModerationComment = sql.NullString{String: *input.ModerationComment, Valid: true}
	}

	review.UpdatedAt = time.Now()

	if err := h.repo.Reviews.Update(c, &review); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении отзыва", err)
		return
	}

	if review.Status == "approved" {
		if err := h.repo.Companies.UpdateRating(c, review.CompanyID); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении рейтинга компании", err)
			return
		}
	}

	updatedReview, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении обновленного отзыва", err)
		return
	}

	utils.Response(c, http.StatusOK, updatedReview)
}

// @Summary Удаление отзыва
// @Description Удаляет отзыв по ID (для администратора)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID отзыва"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/reviews/{id} [delete]
func (h *AdminHandler) DeleteReview(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	review, err := h.repo.Reviews.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Отзыв не найден", err)
		return
	}

	companyID := review.Review.CompanyID

	if err := h.repo.Reviews.Delete(c, id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении отзыва", err)
		return
	}

	if err := h.repo.Companies.UpdateRating(c, companyID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении рейтинга компании", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Отзыв успешно удален"})
}

// @Summary Создание города
// @Description Создает новый город
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CityInput true "Данные города"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/cities [post]
func (h *AdminHandler) CreateCity(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.CityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingCities, _, err := h.repo.Cities.GetAll(c, models.CityFilter{Search: input.Name})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке существования города", err)
		return
	}

	for _, city := range existingCities {
		if city.Name == input.Name && city.Country == input.Country {
			utils.ErrorResponse(c, http.StatusBadRequest, "Город с таким названием уже существует в указанной стране", nil)
			return
		}
	}

	city := &models.City{
		Name:    input.Name,
		Region:  input.Region,
		Country: input.Country,
	}

	id, err := h.repo.Cities.Create(c, city)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании города", err)
		return
	}

	city.ID = id

	utils.Response(c, http.StatusCreated, city)
}

// @Summary Обновление города
// @Description Обновляет существующий город
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID города"
// @Param input body models.CityInput true "Данные города"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/cities/{id} [put]
func (h *AdminHandler) UpdateCity(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.CityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	city, err := h.repo.Cities.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Город не найден", err)
		return
	}

	if city.Name != input.Name || city.Country != input.Country {
		existingCities, _, err := h.repo.Cities.GetAll(c, models.CityFilter{Search: input.Name})
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке существования города", err)
			return
		}

		for _, existingCity := range existingCities {
			if existingCity.Name == input.Name && existingCity.Country == input.Country && existingCity.ID != id {
				utils.ErrorResponse(c, http.StatusBadRequest, "Город с таким названием уже существует в указанной стране", nil)
				return
			}
		}
	}

	city.Name = input.Name
	city.Region = input.Region
	city.Country = input.Country

	if err := h.repo.Cities.Update(c, city); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении города", err)
		return
	}

	utils.Response(c, http.StatusOK, city)
}

// @Summary Удаление города
// @Description Удаляет город по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID города"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/cities/{id} [delete]
func (h *AdminHandler) DeleteCity(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.Cities.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Город не найден", err)
		return
	}

	if err := h.repo.Cities.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении города", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Город успешно удален"})
}

// @Summary Создание индустрии
// @Description Создает новую индустрию
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.IndustryInput true "Данные индустрии"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/industries [post]
func (h *AdminHandler) CreateIndustry(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.IndustryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingIndustry, err := h.repo.Industries.GetByName(c, input.Name)
	if err == nil && existingIndustry != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Индустрия с таким названием уже существует", nil)
		return
	}

	industry := &models.Industry{
		Name:  input.Name,
		Color: input.Color,
	}

	id, err := h.repo.Industries.Create(c, industry)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании индустрии", err)
		return
	}

	industry.ID = id

	utils.Response(c, http.StatusCreated, industry)
}

// @Summary Обновление индустрии
// @Description Обновляет существующую индустрию
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID индустрии"
// @Param input body models.IndustryInput true "Данные индустрии"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/industries/{id} [put]
func (h *AdminHandler) UpdateIndustry(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.IndustryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	industry, err := h.repo.Industries.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Индустрия не найдена", err)
		return
	}

	if industry.Name != input.Name {
		existingIndustry, err := h.repo.Industries.GetByName(c, input.Name)
		if err == nil && existingIndustry != nil && existingIndustry.ID != id {
			utils.ErrorResponse(c, http.StatusBadRequest, "Индустрия с таким названием уже существует", nil)
			return
		}
	}

	industry.Name = input.Name
	if input.Color != "" {
		industry.Color = input.Color
	}

	if err := h.repo.Industries.Update(c, industry); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении индустрии", err)
		return
	}

	utils.Response(c, http.StatusOK, industry)
}

// @Summary Удаление индустрии
// @Description Удаляет индустрию по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID индустрии"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/industries/{id} [delete]
func (h *AdminHandler) DeleteIndustry(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.Industries.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Индустрия не найдена", err)
		return
	}

	if err := h.repo.Industries.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении индустрии", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Индустрия успешно удалена"})
}

// @Summary Создание типа бенефита
// @Description Создает новый тип бенефита
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.BenefitTypeInput true "Данные типа бенефита"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/benefit-types [post]
func (h *AdminHandler) CreateBenefitType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.BenefitTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingType, err := h.repo.BenefitTypes.GetByName(c, input.Name)
	if err == nil && existingType != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Тип бенефита с таким названием уже существует", nil)
		return
	}

	benefitType := &models.BenefitType{
		Name:        input.Name,
		Description: input.Description,
	}

	id, err := h.repo.BenefitTypes.Create(c, benefitType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании типа бенефита", err)
		return
	}

	benefitType.ID = id

	utils.Response(c, http.StatusCreated, benefitType)
}

// @Summary Обновление типа бенефита
// @Description Обновляет существующий тип бенефита
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID типа бенефита"
// @Param input body models.BenefitTypeInput true "Данные типа бенефита"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/benefit-types/{id} [put]
func (h *AdminHandler) UpdateBenefitType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.BenefitTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	benefitType, err := h.repo.BenefitTypes.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Тип бенефита не найден", err)
		return
	}

	if benefitType.Name != input.Name {
		existingType, err := h.repo.BenefitTypes.GetByName(c, input.Name)
		if err == nil && existingType != nil && existingType.ID != id {
			utils.ErrorResponse(c, http.StatusBadRequest, "Тип бенефита с таким названием уже существует", nil)
			return
		}
	}

	benefitType.Name = input.Name
	benefitType.Description = input.Description

	if err := h.repo.BenefitTypes.Update(c, benefitType); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении типа бенефита", err)
		return
	}

	utils.Response(c, http.StatusOK, benefitType)
}

// @Summary Удаление типа бенефита
// @Description Удаляет тип бенефита по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID типа бенефита"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/benefit-types/{id} [delete]
func (h *AdminHandler) DeleteBenefitType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.BenefitTypes.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Тип бенефита не найден", err)
		return
	}

	if err := h.repo.BenefitTypes.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении типа бенефита", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Тип бенефита успешно удален"})
}

// @Summary Создание периода работы
// @Description Создает новый период работы
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.EmploymentPeriodInput true "Данные периода работы"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-periods [post]
func (h *AdminHandler) CreateEmploymentPeriod(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.EmploymentPeriodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingPeriod, err := h.repo.EmploymentPeriods.GetByName(c, input.Name)
	if err == nil && existingPeriod != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Период работы с таким названием уже существует", nil)
		return
	}

	period := &models.EmploymentPeriod{
		Name:        input.Name,
		Description: input.Description,
	}

	id, err := h.repo.EmploymentPeriods.Create(c, period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании периода работы", err)
		return
	}

	period.ID = id

	utils.Response(c, http.StatusCreated, period)
}

// @Summary Обновление периода работы
// @Description Обновляет существующий период работы
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID периода работы"
// @Param input body models.EmploymentPeriodInput true "Данные периода работы"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-periods/{id} [put]
func (h *AdminHandler) UpdateEmploymentPeriod(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.EmploymentPeriodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	period, err := h.repo.EmploymentPeriods.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Период работы не найден", err)
		return
	}

	if period.Name != input.Name {
		existingPeriod, err := h.repo.EmploymentPeriods.GetByName(c, input.Name)
		if err == nil && existingPeriod != nil && existingPeriod.ID != id {
			utils.ErrorResponse(c, http.StatusBadRequest, "Период работы с таким названием уже существует", nil)
			return
		}
	}

	period.Name = input.Name
	period.Description = input.Description

	if err := h.repo.EmploymentPeriods.Update(c, period); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении периода работы", err)
		return
	}

	utils.Response(c, http.StatusOK, period)
}

// @Summary Удаление периода работы
// @Description Удаляет период работы по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID периода работы"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-periods/{id} [delete]
func (h *AdminHandler) DeleteEmploymentPeriod(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.EmploymentPeriods.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Период работы не найден", err)
		return
	}

	if err := h.repo.EmploymentPeriods.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении периода работы", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Период работы успешно удален"})
}

// @Summary Создание типа занятости
// @Description Создает новый тип занятости
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.EmploymentTypeInput true "Данные типа занятости"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-types [post]
func (h *AdminHandler) CreateEmploymentType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	var input models.EmploymentTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingType, err := h.repo.EmploymentTypes.GetByName(c, input.Name)
	if err == nil && existingType != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Тип занятости с таким названием уже существует", nil)
		return
	}

	empType := &models.EmploymentType{
		Name:        input.Name,
		Description: input.Description,
	}

	id, err := h.repo.EmploymentTypes.Create(c, empType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании типа занятости", err)
		return
	}

	empType.ID = id

	utils.Response(c, http.StatusCreated, empType)
}

// @Summary Обновление типа занятости
// @Description Обновляет существующий тип занятости
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID типа занятости"
// @Param input body models.EmploymentTypeInput true "Данные типа занятости"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-types/{id} [put]
func (h *AdminHandler) UpdateEmploymentType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	var input models.EmploymentTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	empType, err := h.repo.EmploymentTypes.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Тип занятости не найден", err)
		return
	}

	if empType.Name != input.Name {
		existingType, err := h.repo.EmploymentTypes.GetByName(c, input.Name)
		if err == nil && existingType != nil && existingType.ID != id {
			utils.ErrorResponse(c, http.StatusBadRequest, "Тип занятости с таким названием уже существует", nil)
			return
		}
	}

	empType.Name = input.Name
	empType.Description = input.Description

	if err := h.repo.EmploymentTypes.Update(c, empType); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении типа занятости", err)
		return
	}

	utils.Response(c, http.StatusOK, empType)
}

// @Summary Удаление типа занятости
// @Description Удаляет тип занятости по ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID типа занятости"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 403 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /admin/employment-types/{id} [delete]
func (h *AdminHandler) DeleteEmploymentType(c *gin.Context) {
	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists || roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return
	}

	_, err = h.repo.EmploymentTypes.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Тип занятости не найден", err)
		return
	}

	if err := h.repo.EmploymentTypes.Delete(c, id); err != nil {
		if strings.Contains(err.Error(), "используется") {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении типа занятости", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, gin.H{"message": "Тип занятости успешно удален"})
}
