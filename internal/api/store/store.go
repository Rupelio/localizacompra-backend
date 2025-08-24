package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	CNPJ      string    `json:"cnpj"`
}

type CreateStoreRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	CNPJ    string `json:"cnpj"`
}

type AdminDataRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type StoreWithAdminRequest struct {
	StoreName    string           `json:"store_name"`
	StoreAddress string           `json:"store_address"`
	CNPJ         string           `json:"cpnj"`
	Admin        AdminDataRequest `json:"admin"`
}

type Repository interface {
	Create(ctx context.Context, store Store) (Store, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, store Store) (Store, error)
	GetAll(ctx context.Context) ([]Store, error)
}

type Service interface {
	Create(ctx context.Context, store Store) (Store, error)
	CreateStoreWithAdmin(ctx context.Context, req StoreWithAdminRequest) (Store, error)
	GetAll(ctx context.Context) ([]Store, error)
}
