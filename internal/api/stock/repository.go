package stock

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

func (r *pgxRepository) Create(ctx context.Context, item StockItem) (StockItem, error) {
	query := `INSERT INTO stock_items (store_id, product_id, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRow(ctx, query,
		item.StoreID,
		item.ProductID,
		item.Price,
		item.Quantity,
	).Scan(&item.ID)

	if err != nil {
		return StockItem{}, err
	}

	return item, nil
}

func (r *pgxRepository) GetAllByStoreId(ctx context.Context, storeID int64) ([]ProductStockDetail, error) {
	query := `SELECT
				products.id,
				products.name,
				products.description,
				stock_items.price,
				stock_items.quantity
			FROM stock_items
			JOIN products ON stock_items.product_id = products.id
			WHERE stock_items.store_id = $1;`

	rows, err := r.db.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var productStockDetail []ProductStockDetail

	for rows.Next() {
		var p ProductStockDetail

		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}

		productStockDetail = append(productStockDetail, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productStockDetail, nil
}
