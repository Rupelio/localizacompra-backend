package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(dbpool *pgxpool.Pool) Repository {
	return &pgxRepository{
		db: dbpool,
	}
}

func (r *pgxRepository) Create(ctx context.Context, store Store) (Store, error) {
	query := `INSERT INTO stores (name, address) VALUES ($1, $2) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, store.Name, store.Address).Scan(&store.ID, &store.CreatedAt)
	if err != nil {
		return Store{}, err
	}

	return store, nil
}

func (r *pgxRepository) GetAll(ctx context.Context) ([]Store, error) {
	query := `SELECT id, name, address, created_at FROM stores`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stores []Store

	for rows.Next() {
		var s Store
		err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.CreatedAt)

		if err != nil {
			return nil, err
		}

		stores = append(stores, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}
