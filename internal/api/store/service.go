package store

import (
	"context"
	"errors"
	"localiza-compra/backend/internal/api/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type storeService struct {
	db       *pgxpool.Pool
	repo     Repository
	userRepo user.Repository
}

func NewService(db *pgxpool.Pool, r Repository, ur user.Repository) Service {
	return &storeService{
		db:       db,
		repo:     r,
		userRepo: ur,
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

func (s *storeService) CreateStoreWithAdmin(ctx context.Context, req StoreWithAdminRequest) (Store, error) {
	if req.StoreName == "" {
		return Store{}, errors.New("o nome da loja não pode ser vazia")
	}
	if req.StoreAddress == "" {
		return Store{}, errors.New("o endereço da loja não pode ser vazio")
	}
	if req.CNPJ == "" {
		return Store{}, errors.New("o CNPJ da loja não pode ser vazio")
	}
	if len(req.CNPJ) != 14 {
		return Store{}, errors.New("o CNPJ da loja precisa ter 14 caracteres")
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Store{}, err
	}

	defer tx.Rollback(ctx)

	storeToCreate := Store{
		Name:    req.StoreName,
		Address: req.StoreAddress,
		CNPJ:    req.CNPJ,
	}

	createdStore, err := s.repo.CreateWithTx(ctx, tx, storeToCreate)
	if err != nil {
		return Store{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return Store{}, err
	}

	adminToCreate := user.User{
		Name:         req.Admin.Name,
		Email:        req.Admin.Email,
		PasswordHash: string(hashedPassword),
		Phone:        req.Admin.Phone,
		Role:         "store_admin",
		StoreID:      &createdStore.ID,
	}

	_, err = s.userRepo.CreateWithTx(ctx, tx, adminToCreate)
	if err != nil {
		return Store{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return Store{}, err
	}

	return createdStore, nil
}
