package shoppinglist

import (
	"context"
	"errors"
	"time"
)

var ErrShoppingListNotFound = errors.New("lista não encontrada")
var ErrShoppingListItemNotFound = errors.New("produto não encontrado")

type ShoppingList struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ItemCount int       `json:"item_count"`
}

type ShoppingListItem struct {
	ID             int64 `json:"id"`
	ShoppingListID int64 `json:"shopping_list_id"`
	ProductID      int64 `json:"product_id"`
	Quantity       int   `json:"quantity"`
	IsChecked      bool  `json:"is_checked"`
}

type CreateShoppingListItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CreateShoppingListRequest struct {
	Name string `json:"name"`
}

type ListItemDetail struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"product_name"`
	Quantity  int    `json:"quantity"`
	IsChecked bool   `json:"is_checked"`
}

type UpdateItemRequest struct {
	IsChecked bool `json:"is_checked"`
}

type OptimizedListItem struct {
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	IsChecked   bool    `json:"is_checked"`
	Price       float64 `json:"price"`
	Sector      string  `json:"sector"`
}

type Repository interface {
	CreateList(ctx context.Context, list ShoppingList) (ShoppingList, error)
	CreateItem(ctx context.Context, item ShoppingListItem) (ShoppingListItem, error)
	GetShoppingListByID(ctx context.Context, id int64) (ShoppingList, error)
	GetAllByUserID(ctx context.Context, userID int64) ([]ShoppingList, error)
	GetAllItemsByListID(ctx context.Context, listID int64) ([]ListItemDetail, error)
	UpdateItemStatus(ctx context.Context, itemID int64, isChecked bool) error
	GetOptimizedList(ctx context.Context, listID int64, storeID int64) ([]OptimizedListItem, error)
}

type Service interface {
	CreateList(ctx context.Context, list ShoppingList) (ShoppingList, error)
	CreateItem(ctx context.Context, userID int64, item ShoppingListItem) (ShoppingListItem, error)
	GetAllByUserID(ctx context.Context, userID int64) ([]ShoppingList, error)
	GetAllItemsByListID(ctx context.Context, userID, listID int64) ([]ListItemDetail, error)
	UpdateItemStatus(ctx context.Context, userID, listID, itemID int64, isChecked bool) error
	GetOptimizedList(ctx context.Context, userID, listID, storeID int64) ([]OptimizedListItem, error)
}
