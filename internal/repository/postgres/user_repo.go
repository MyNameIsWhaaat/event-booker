package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
		SELECT id, email, created_at
		FROM users
		WHERE email = $1
	`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, email string) (domain.User, error) {
	const q = `
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id, email, created_at
	`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) GetOrCreateByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
		INSERT INTO users (email)
		VALUES ($1)
		ON CONFLICT (email) DO UPDATE
		SET email = EXCLUDED.email
		RETURNING id, email, created_at
	`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}