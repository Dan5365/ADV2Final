package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, name, email, password string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
		 RETURNING id, name, email, password, created_at`,
		name, email, password,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password, created_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, id int64, name, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`UPDATE users SET name=$1, email=$2 WHERE id=$3
		 RETURNING id, name, email, password, created_at`,
		name, email, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}
