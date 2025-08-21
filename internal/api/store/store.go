package store

import (
	"context"
	"time"
)

type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStoreRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Repository interface {
	Create(ctx context.Context, store Store) (Store, error)
	GetAll(ctx context.Context) ([]Store, error)
}

type Service interface {
	Create(ctx context.Context, store Store) (Store, error)
	GetAll(ctx context.Context) ([]Store, error)
}
