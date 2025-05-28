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
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewUserHandler(postgres *db.PostgreSQL, cfg *config.Config) *UserHandler {
	repo := repository.NewRepository(postgres)
	return &UserHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// @Summary Профиль пользователя
// @Description Возвращает профиль текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.ResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	user, err := h.repo.Users.GetByID(c, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Пользователь не найден", nil)
		return
	}

	utils.Response(c, http.StatusOK, user.ToProfile())
}

// @Summary Обновление профиля
// @Description Обновляет профиль текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.UserUpdateInput true "Данные для обновления"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	var input models.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	user, err := h.repo.Users.GetByID(c, userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Пользователь не найден", nil)
		return
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при хешировании пароля", err)
			return
		}
		user.PasswordHash = string(passwordHash)
	}

	user.UpdatedAt = time.Now()

	err = h.repo.Users.Update(c, user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении профиля", err)
		return
	}

	utils.Response(c, http.StatusOK, user.ToProfile())
}

// @Summary Отзывы пользователя
// @Description Возвращает отзывы текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Статус отзывов (pending, approved, rejected)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /users/me/reviews [get]
func (h *UserHandler) GetUserReviews(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	var filter models.ReviewFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации параметров", err)
		return
	}

	userId := userID.(int)
	filter.UserID = &userId

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	reviews, total, err := h.repo.Reviews.GetByUser(c, userId, filter)
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
