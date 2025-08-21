package user

import (
	"context"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("usuario não encontrado")

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Phone        string    `json:"phone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"-"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	// Update(ctx context.Context, user User) (User, error)
	// Delete(ctx context.Context, id int64) error
}

type Service interface {
	GetByID(ctx context.Context, id int64) (User, error)
	Create(ctx context.Context, user User) (User, error)
	// Update(ctx context.Context, user User) (User, error)
	// Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (string, error)
}
