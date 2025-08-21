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
	query := `SELECT id, name, email, password_hash, phone, created_at, role FROM users WHERE email = $1`

	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.CreatedAt,
		&user.Role,
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
	query := `SELECT id, name, email, password_hash, phone, created_at, role FROM users WHERE id = $1`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.CreatedAt,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}
