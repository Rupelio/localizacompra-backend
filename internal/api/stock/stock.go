package stock

import (
	"context"
	"errors"
)

var ErrStockItemNotFount = errors.New("item de estoque não encontrado")

type StockItem struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"id_product"`
	StoreID   int64   `json:"id_store"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type CreateStockItemRequest struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ProductStockDetail struct {
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Repository interface {
	Create(ctx context.Context, item StockItem) (StockItem, error)
	GetAllByStoreId(ctx context.Context, storeID int64) ([]ProductStockDetail, error)
}

type Service interface {
	Create(ctx context.Context, item StockItem) (StockItem, error)
	GetAllByStoreId(ctx context.Context, storeID int64) ([]ProductStockDetail, error)
}
