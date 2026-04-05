package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/model"
)

type usersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(ctx context.Context, user *model.User) error {
	const q = `
		INSERT INTO users (email, password_hash, name, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, q, user.Email, user.PasswordHash, user.Name).Scan(&user.ID, &user.CreatedAt)
}

func (r *usersRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const q = `
		SELECT id, email, password_hash, name, created_at
		FROM users
		WHERE email = $1`
	row := r.db.QueryRow(ctx, q, email)
	return scanUser(row)
}

func (r *usersRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	const q = `
		SELECT id, email, password_hash, name, created_at
		FROM users
		WHERE id = $1`
	row := r.db.QueryRow(ctx, q, id)
	return scanUser(row)
}

func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
