package store

import (
	"context"
	"errors"
)

type storeService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &storeService{
		repo: r,
	}
}

func (s *storeService) Create(ctx context.Context, store Store) (Store, error) {
	if store.Name == "" {
		return Store{}, errors.New("o nome da loja não pode ser vazia")
	}
	if store.Address == "" {
		return Store{}, errors.New("o endereço da loja não pode ser vazio")
	}
	if store.CNPJ == "" {
		return Store{}, errors.New("o CNPJ da loja não pode ser vazio")
	}
	if len(store.CNPJ) != 14 {
		return Store{}, errors.New("o CNPJ da loja precisa ter 14 caracteres")
	}
	return s.repo.Create(ctx, store)
}

func (s *storeService) GetAll(ctx context.Context) ([]Store, error) {
	return s.repo.GetAll(ctx)
}
