package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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
	StoreID      *int64    `json:"store_id"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRole struct {
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

type Repository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, user User) (User, error)
	// Update(ctx context.Context, user User) (User, error)
	// Delete(ctx context.Context, id int64) error
	UpdateRole(ctx context.Context, email string, role string) error
}

type Service interface {
	GetByID(ctx context.Context, id int64) (User, error)
	Create(ctx context.Context, user User) (User, error)
	// Update(ctx context.Context, user User) (User, error)
	// Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (string, error)
}
