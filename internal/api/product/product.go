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
	Brand       string    `json:"brand"`
	ImageUrl    string    `json:"image_url"`
	CategoryID  *int64    `json:"category_id,omitempty"`
}

// UpdateProductRequest é o DTO para atualizações parciais (PATCH)
type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	ImageUrl    *string `json:"image_url,omitempty"`
	CategoryID  *int64  `json:"category_id,omitempty"`
}

type Service interface {
	GetAll(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id int64) error
	SearchByName(ctx context.Context, name string) ([]Product, error)
	PartialUpdate(ctx context.Context, id int64, req UpdateProductRequest) error
}

type Repository interface {
	GetAll(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id int64) error
	SearchByName(ctx context.Context, name string) ([]Product, error)
	PartialUpdate(ctx context.Context, id int64, req UpdateProductRequest) error
}
