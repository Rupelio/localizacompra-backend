package stock

import "context"

type stockItemService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &stockItemService{
		repo: r,
	}
}

func (s *stockItemService) Create(ctx context.Context, item StockItem) (StockItem, error) {
	return s.repo.Create(ctx, item)
}

func (s *stockItemService) GetAllByStoreId(ctx context.Context, storeID int64) ([]ProductStockDetail, error) {
	return s.repo.GetAllByStoreId(ctx, storeID)
}
