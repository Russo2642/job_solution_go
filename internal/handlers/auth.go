package handlers

import (
	"net/http"
	"time"

	"job_solition/internal/config"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo repository.Repository
	jwt  *utils.JWT
}

func NewAuthHandler(repo *repository.Repository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		repo: *repo,
		jwt: utils.NewJWT(
			cfg.JWT.Secret,
			cfg.JWT.ExpiresIn,
			cfg.JWT.RefreshExpiresIn,
		),
	}
}

// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserRegisterInput true "Данные для регистрации"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 409 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.UserRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	if input.Password != input.PasswordConfirm {
		utils.ErrorResponse(c, http.StatusBadRequest, "Пароли не совпадают", nil)
		return
	}

	existingUser, err := h.repo.Users.GetByEmail(c, input.Email)
	if err == nil && existingUser != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Пользователь с таким email уже существует", nil)
		return
	} else if err != nil && err.Error() != "пользователь не найден" {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке существующего пользователя", err)
		return
	}

	user, err := models.NewUser(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании пользователя", err)
		return
	}

	userID, err := h.repo.Users.Create(c, user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении пользователя", err)
		return
	}

	user.ID = userID

	token, err := h.jwt.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании токена", err)
		return
	}

	refreshToken := models.NewRefreshToken(user.ID, h.jwt.RefreshExpiresIn)
	_, err = h.repo.RefreshTokens.Create(c, &refreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании refresh токена", err)
		return
	}

	utils.Response(c, http.StatusCreated, gin.H{
		"user": user.ToProfile(),
		"tokens": gin.H{
			"access_token":  token,
			"refresh_token": refreshToken.Token,
		},
	})
}

// @Summary Вход в систему
// @Description Аутентифицирует пользователя и выдает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserLoginInput true "Данные для входа"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	user, err := h.repo.Users.GetByEmail(c, input.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Неверный email или пароль", nil)
		return
	}

	if !user.ComparePassword(input.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Неверный email или пароль", nil)
		return
	}

	token, err := h.jwt.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании токена", err)
		return
	}

	refreshToken := models.NewRefreshToken(user.ID, h.jwt.RefreshExpiresIn)
	_, err = h.repo.RefreshTokens.Create(c, &refreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании refresh токена", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"user": user.ToProfile(),
		"tokens": gin.H{
			"access_token":  token,
			"refresh_token": refreshToken.Token,
		},
	})
}

// @Summary Обновление токенов
// @Description Обновляет пару access/refresh токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Refresh токен"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	refreshToken, err := h.repo.RefreshTokens.GetByToken(c, input.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Недействительный refresh токен", nil)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Refresh токен просрочен", nil)
		return
	}

	user, err := h.repo.Users.GetByID(c, refreshToken.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Пользователь не найден", nil)
		return
	}

	err = h.repo.RefreshTokens.DeleteByToken(c, input.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении старого refresh токена", err)
		return
	}

	token, err := h.jwt.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании токена", err)
		return
	}

	newRefreshToken := models.NewRefreshToken(user.ID, h.jwt.RefreshExpiresIn)
	_, err = h.repo.RefreshTokens.Create(c, &newRefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при создании refresh токена", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"tokens": gin.H{
			"access_token":  token,
			"refresh_token": newRefreshToken.Token,
		},
	})
}

// @Summary Выход из системы
// @Description Выход из системы, удаление refresh токена
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Refresh токен"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	err := h.repo.RefreshTokens.DeleteByToken(c, input.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при выходе из системы", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Успешный выход из системы",
	})
}

// @Summary Восстановление пароля
// @Description Позволяет пользователю сбросить пароль, указав email и новый пароль
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.ForgotPasswordInput true "Email пользователя и новый пароль"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input models.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	if input.Password != input.PasswordConfirm {
		utils.ErrorResponse(c, http.StatusBadRequest, "Пароли не совпадают", nil)
		return
	}

	user, err := h.repo.Users.GetByEmail(c, input.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Пользователь с указанным email не найден", nil)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при хешировании пароля", err)
		return
	}

	user.PasswordHash = string(passwordHash)
	user.UpdatedAt = time.Now()

	if err := h.repo.Users.Update(c, user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении пароля", err)
		return
	}

	if err := h.repo.RefreshTokens.DeleteByUserID(c, user.ID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при удалении refresh токенов", err)
		return
	}

	utils.Response(c, http.StatusOK, gin.H{
		"message": "Пароль успешно изменен",
	})
}
