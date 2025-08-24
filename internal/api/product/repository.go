package product

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
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
	query := "SELECT id, name, description, created_at, brand, image_url, category_id FROM products"

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.Brand, &p.ImageUrl, &p.CategoryID)
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
		INSERT INTO products (name, description, brand, image_url, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Brand, product.ImageUrl, product.CategoryID).Scan(&product.ID, &product.CreatedAt)

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// Update atualiza um produto já existente no banco de dados
func (r *pgxProductRepository) Update(ctx context.Context, product Product) (Product, error) {
	query := `
		UPDATE products
		SET name = $1, description = $2, brand = $3, image_url = $4
		WHERE id = $5
		RETURNING id, name, description, created_at, brand, image_url
	`

	var updatedProduct Product
	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Brand,
		product.ImageUrl,
		product.ID,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Description,
		&updatedProduct.CreatedAt,
		&updatedProduct.Brand,
		&updatedProduct.ImageUrl,
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

	query := `SELECT id, name, description, created_at, brand, image_url FROM products WHERE name ILIKE $1`

	rows, err := r.db.Query(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.Brand, &p.ImageUrl)
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

func (r *pgxProductRepository) PartialUpdate(ctx context.Context, id int64, req UpdateProductRequest) error {
	updateBuilder := sq.Update("products").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	if req.Name != nil {
		updateBuilder = updateBuilder.Set("name", *req.Name)
	}
	if req.Description != nil {
		updateBuilder = updateBuilder.Set("description", *req.Description)
	}
	if req.Brand != nil {
		updateBuilder = updateBuilder.Set("brand", *req.Brand)
	}
	if req.ImageUrl != nil {
		updateBuilder = updateBuilder.Set("image_url", *req.ImageUrl)
	}
	if req.CategoryID != nil {
		updateBuilder = updateBuilder.Set("category_id", *req.CategoryID)
	}

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}
