package product

import (
	"context"
	"errors"
	"time"
)

var ErrProductNotFound = errors.New("produto não encontrado")

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Service interface {
	GetAll(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id int64) error
	SearchByName(ctx context.Context, name string) ([]Product, error)
}

type Repository interface {
	GetAll(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id int64) error
	SearchByName(ctx context.Context, name string) ([]Product, error)
}
