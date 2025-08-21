package shoppinglist

import (
	"context"
	"errors"

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

func (r *pgxRepository) CreateList(ctx context.Context, list ShoppingList) (ShoppingList, error) {
	query := `INSERT INTO shopping_lists (user_id, name) VALUES ($1, $2) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, list.UserID, list.Name).Scan(&list.ID, &list.CreatedAt)
	if err != nil {
		return ShoppingList{}, err
	}

	return list, nil
}

func (r *pgxRepository) CreateItem(ctx context.Context, item ShoppingListItem) (ShoppingListItem, error) {
	query := `INSERT INTO shopping_list_items (shopping_list_id, product_id, quantity, is_checked)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

	err := r.db.QueryRow(ctx, query, item.ShoppingListID, item.ProductID, item.Quantity, item.IsChecked).Scan(&item.ID)
	if err != nil {
		return ShoppingListItem{}, err
	}

	return item, nil
}

func (r *pgxRepository) GetShoppingListByID(ctx context.Context, id int64) (ShoppingList, error) {
	query := `SELECT id, user_id, name, created_at FROM shopping_lists WHERE id = $1`

	var list ShoppingList

	err := r.db.QueryRow(ctx, query, id).Scan(&list.ID, &list.UserID, &list.Name, &list.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ShoppingList{}, ErrShoppingListNotFound
		}
		return ShoppingList{}, err
	}

	return list, nil
}

func (r *pgxRepository) GetAllByUserID(ctx context.Context, userID int64) ([]ShoppingList, error) {
	query := `SELECT
			sl.id,
			sl.user_id,
			sl.name,
			sl.created_at,
			COUNT(sli.id) AS item_count
		FROM
			shopping_lists sl
		LEFT JOIN
			shopping_list_items sli ON sl.id = sli.shopping_list_id
		WHERE
			sl.user_id = $1
		GROUP BY
			sl.id;
		`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lists []ShoppingList
	for rows.Next() {
		var l ShoppingList

		err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.CreatedAt, &l.ItemCount)
		if err != nil {
			return nil, err
		}

		lists = append(lists, l)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *pgxRepository) GetAllItemsByListID(ctx context.Context, listID int64) ([]ListItemDetail, error) {
	query := `SELECT
			sli.id,
			sli.product_id,
			p.name AS product_name,
			sli.quantity,
			sli.is_checked
		FROM
			shopping_list_items sli
		JOIN
			products p ON sli.product_id = p.id
		WHERE
			sli.shopping_list_id = $1;
		`
	rows, err := r.db.Query(ctx, query, listID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]ListItemDetail, 0)

	for rows.Next() {
		var i ListItemDetail

		err := rows.Scan(&i.ID, &i.ProductID, &i.Name, &i.Quantity, &i.IsChecked)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *pgxRepository) UpdateItemStatus(ctx context.Context, itemID int64, isChecked bool) error {
	query := `UPDATE shopping_list_items SET is_checked = $1 WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, isChecked, itemID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrShoppingListItemNotFound
	}

	return nil
}
