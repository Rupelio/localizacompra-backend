package category

import (
	"context"
	"errors"
)

type categoryService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &categoryService{
		repo: r,
	}
}

func (s *categoryService) GetByID(ctx context.Context, id int64) (Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) Create(ctx context.Context, category Category) (Category, error) {
	// Validação básica
	if category.Name == "" {
		return Category{}, errors.New("o nome da categoria não pode ser vazio")
	}

	// Se for uma sub-categoria (ParentID não é nulo)
	if category.ParentID != nil {
		// Verificamos se a categoria pai existe
		_, err := s.repo.GetByID(ctx, *category.ParentID)
		if err != nil {
			// Se repo.GetByID devolveu um erro, significa que o pai não existe
			return Category{}, errors.New("a categoria pai especificada não existe")
		}
	}

	// Se a categoria pai existe (ou se não há pai), criamos a nova categoria
	return s.repo.Create(ctx, category)
}

func (s *categoryService) PartialUpdate(ctx context.Context, id int64, req UpdateCategoryRequest) error {
	return s.repo.PartialUpdate(ctx, id, req)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *categoryService) GetAll(ctx context.Context) ([]Category, error) {
	return s.repo.GetAll(ctx)
}
