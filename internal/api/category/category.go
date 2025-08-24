package category

import (
	"context"
	"errors"
)

var ErrCategoryNotFound = errors.New("categoria não encontrada")

type Category struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name"`
	ParentID *int64  `json:"parent_id"`
}

type Repository interface {
	Create(ctx context.Context, category Category) (Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int64) (Category, error)
	PartialUpdate(ctx context.Context, id int64, req UpdateCategoryRequest) error
	Delete(ctx context.Context, id int64) error
}

type Service interface {
	Create(ctx context.Context, category Category) (Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int64) (Category, error)
	PartialUpdate(ctx context.Context, id int64, req UpdateCategoryRequest) error
	Delete(ctx context.Context, id int64) error
}
