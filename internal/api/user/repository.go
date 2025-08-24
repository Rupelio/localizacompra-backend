package user

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

func (r *pgxRepository) Create(ctx context.Context, user User) (User, error) {
	query := `
			INSERT INTO users (name, email, password_hash, phone)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, user.Phone).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *pgxRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, name, email, password_hash, phone, created_at, role, store_id FROM users WHERE email = $1`

	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.CreatedAt,
		&user.Role,
		&user.StoreID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (r *pgxRepository) GetByID(ctx context.Context, id int64) (User, error) {
	query := `SELECT id, name, email, password_hash, phone, created_at, role, store_id FROM users WHERE id = $1`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.CreatedAt,
		&user.Role,
		&user.StoreID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (r *pgxRepository) UpdateRole(ctx context.Context, email string, role string) error {
	query := `UPDATE users SET role = $1 WHERE email = $2`

	tag, err := r.db.Exec(ctx, query, role, email)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *pgxRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, user User) (User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, phone, role, store_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := tx.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.Role,
		user.StoreID,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
