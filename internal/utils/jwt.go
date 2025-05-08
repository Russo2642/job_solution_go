package utils

import (
	"fmt"
	"time"

	"job_solition/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

type AuthClaims struct {
	UserID int             `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewJWT(secret string, expiresIn, refreshExpiresIn time.Duration) *JWT {
	return &JWT{
		Secret:           secret,
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}
}

func (j *JWT) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи токена: %w", err)
	}

	return tokenString, nil
}

func (j *JWT) ValidateToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка валидации токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("недействительный токен")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, fmt.Errorf("неверный формат токена")
	}

	return claims, nil
}
