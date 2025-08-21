package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxProductRepository é a implementação do ProductRepository usando pgxpool.
type pgxProductRepository struct {
	db *pgxpool.Pool
}

// NewPgxProductRepository cria uma nova instância do nosso repositório.
func NewRepository(dbpool *pgxpool.Pool) Repository {
	return &pgxProductRepository{
		db: dbpool,
	}
}

// GetAll busca todos os produtos no banco de dados.
func (r *pgxProductRepository) GetAll(ctx context.Context) ([]Product, error) {
	query := "SELECT id, name, description, created_at FROM products"

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}

// Create insere um novo produto no banco de dados.
func (r *pgxProductRepository) Create(ctx context.Context, product Product) (Product, error) {

	query := `
		INSERT INTO products (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, product.Name, product.Description).Scan(&product.ID, &product.CreatedAt)

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// Update atualiza um produto já existente no banco de dados
func (r *pgxProductRepository) Update(ctx context.Context, product Product) (Product, error) {
	query := `
		UPDATE products
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at
	`

	var updatedProduct Product
	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.ID,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Description,
		&updatedProduct.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, ErrProductNotFound
		}
		return Product{}, err
	}

	return product, nil
}

// Delete remove um produto do banco de dados pelo seu ID
func (r *pgxProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r *pgxProductRepository) SearchByName(ctx context.Context, name string) ([]Product, error) {
	searchTerm := "%" + name + "%"

	query := `SELECT id, name, description, created_at FROM products WHERE name ILIKE $1`

	rows, err := r.db.Query(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
