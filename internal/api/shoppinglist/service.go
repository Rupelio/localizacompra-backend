package shoppinglist

import (
	"context"
	"errors"
)

type shoppingService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &shoppingService{
		repo: r,
	}
}

func (s *shoppingService) CreateItem(ctx context.Context, userID int64, item ShoppingListItem) (ShoppingListItem, error) {
	list, err := s.repo.GetShoppingListByID(ctx, item.ShoppingListID)
	if err != nil {
		return ShoppingListItem{}, err
	}

	if list.UserID != userID {
		return ShoppingListItem{}, errors.New("não autorizado: você não é o dono desta lista")
	}
	return s.repo.CreateItem(ctx, item)
}

func (s *shoppingService) CreateList(ctx context.Context, list ShoppingList) (ShoppingList, error) {
	return s.repo.CreateList(ctx, list)
}

func (s *shoppingService) GetAllByUserID(ctx context.Context, userID int64) ([]ShoppingList, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *shoppingService) GetAllItemsByListID(ctx context.Context, userID, listID int64) ([]ListItemDetail, error) {
	list, err := s.repo.GetShoppingListByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	// 2. A VERIFICAÇÃO DE SEGURANÇA!
	if list.UserID != userID {
		return nil, errors.New("não autorizado")
	}

	return s.repo.GetAllItemsByListID(ctx, listID)
}

func (s *shoppingService) UpdateItemStatus(ctx context.Context, userID, listID, itemID int64, isChecked bool) error {
	// 1. Buscamos a lista para descobrir o dono
	list, err := s.repo.GetShoppingListByID(ctx, listID)
	if err != nil {
		return err // Devolve o erro de "não encontrado" se a lista não existir
	}

	// 2. A verificação de segurança!
	if list.UserID != userID {
		return errors.New("não autorizado")
	}

	// 3. Se a verificação passar, mandamos o repositório fazer a atualização.
	// Note que o repositório só precisa do itemID e do isChecked.
	return s.repo.UpdateItemStatus(ctx, itemID, isChecked)
}

func (s *shoppingService) GetOptimizedList(ctx context.Context, userID, listID, storeID int64) ([]OptimizedListItem, error) {
	// Verificação de segurança: o utilizador é o dono da lista?
	list, err := s.repo.GetShoppingListByID(ctx, listID)
	if err != nil {
		return nil, err
	}
	if list.UserID != userID {
		return nil, errors.New("não autorizado")
	}

	// Se for, busca a lista otimizada
	return s.repo.GetOptimizedList(ctx, listID, storeID)
}
