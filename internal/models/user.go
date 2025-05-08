package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleUser UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name,omitempty" db:"first_name"`
	LastName     string    `json:"last_name,omitempty" db:"last_name"`
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserProfile struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterInput struct {
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"omitempty,phone"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"omitempty"`
	LastName        string `json:"last_name" binding:"omitempty"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateInput struct {
	Phone     *string `json:"phone" binding:"omitempty,phone"`
	FirstName *string `json:"first_name" binding:"omitempty"`
	LastName  *string `json:"last_name" binding:"omitempty"`
	Password  *string `json:"password" binding:"omitempty,min=8"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

type PasswordResetToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUser(input UserRegisterInput) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(passwordHash),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Role:         RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}

func NewRefreshToken(userID int, expiresIn time.Duration) RefreshToken {
	now := time.Now()
	return RefreshToken{
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: now.Add(expiresIn),
		CreatedAt: now,
	}
}

func NewPasswordResetToken(userID int) PasswordResetToken {
	now := time.Now()
	return PasswordResetToken{
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
}
