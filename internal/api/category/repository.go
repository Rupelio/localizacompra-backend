package category

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (r *pgxRepository) Create(ctx context.Context, category Category) (Category, error) {
	query := `INSERT INTO categories (name, parent_id)
			VALUES ($1, $2)
			RETURNING id`

	err := r.db.QueryRow(ctx, query, category.Name, category.ParentID).Scan(&category.ID)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (r *pgxRepository) GetByID(ctx context.Context, id int64) (Category, error) {
	query := `SELECT id, name, parent_id FROM categories WHERE id = $1`

	var c Category

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.ParentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Category{}, ErrCategoryNotFound
		}
		return Category{}, err
	}

	return c, nil
}

func (r *pgxRepository) PartialUpdate(ctx context.Context, id int64, req UpdateCategoryRequest) error {
	updateBuilder := sq.Update("categories").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)
	if req.Name != nil {
		updateBuilder = updateBuilder.Set("name", *req.Name)
	}
	if req.ParentID != nil {
		updateBuilder = updateBuilder.Set("parent_id", *req.ParentID)
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
		return ErrCategoryNotFound
	}

	return nil
}

func (r *pgxRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *pgxRepository) GetAll(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name, parent_id FROM categories`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category

		err := rows.Scan(&c.ID, &c.Name, &c.ParentID)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
