package product

import (
	"context"
	"errors"
)

type productService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &productService{
		repo: r,
	}
}

func (s *productService) GetAll(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) Create(ctx context.Context, product Product) (Product, error) {
	if product.Name == "" {
		return Product{}, errors.New("o nome do produto não pode ser vazio")
	}

	return s.repo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, product Product) (Product, error) {
	if product.Name == "" {
		return Product{}, errors.New("o nome do produto não pode ser vazio")
	}

	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) SearchByName(ctx context.Context, name string) ([]Product, error) {
	return s.repo.SearchByName(ctx, name)
}

func (s *productService) PartialUpdate(ctx context.Context, id int64, req UpdateProductRequest) error {
	return s.repo.PartialUpdate(ctx, id, req)
}
