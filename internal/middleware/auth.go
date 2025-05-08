package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"job_solition/internal/config"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey          = "user_id"
	UserKey            = "user"
	RoleKey            = "role"
	IsAuthenticatedKey = "is_authenticated"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	jwtUtil := utils.NewJWT(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshExpiresIn)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Отсутствует заголовок Authorization", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Неверный формат заголовка Authorization", nil)
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Недействительный токен", err)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Set(IsAuthenticatedKey, true)

		c.Next()
	}
}

func OptionalAuth(cfg *config.Config) gin.HandlerFunc {
	jwtUtil := utils.NewJWT(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshExpiresIn)

	return func(c *gin.Context) {
		c.Set(IsAuthenticatedKey, false)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Printf("OptionalAuth: заголовок Authorization отсутствует\n")
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("OptionalAuth: неверный формат заголовка: %s\n", authHeader)
			c.Next()
			return
		}

		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			fmt.Printf("OptionalAuth: ошибка валидации токена: %v\n", err)
			c.Next()
			return
		}

		fmt.Printf("OptionalAuth: пользователь аутентифицирован, ID: %d, роль: %s\n", claims.UserID, claims.Role)

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Set(IsAuthenticatedKey, true)

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthenticated, exists := c.Get(IsAuthenticatedKey)

		fmt.Printf("RequireAuth: exists=%v, isAuthenticated=%v\n", exists, isAuthenticated)

		authHeader := c.GetHeader("Authorization")
		fmt.Printf("RequireAuth: Authorization header=%s\n", authHeader)

		if !exists || isAuthenticated != true {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRoleMiddleware(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(RoleKey)
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
			c.Abort()
			return
		}

		role, ok := roleValue.(models.UserRole)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Неверный формат роли", nil)
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if role == r {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func LoadUserMiddleware(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get(UserIDKey)
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
			c.Abort()
			return
		}

		userID, ok := userIDValue.(int)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Неверный формат ID пользователя", nil)
			c.Abort()
			return
		}

		user, err := repo.Users.GetByID(c, userID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Пользователь не найден", nil)
			c.Abort()
			return
		}

		c.Set(UserKey, user)

		c.Next()
	}
}
