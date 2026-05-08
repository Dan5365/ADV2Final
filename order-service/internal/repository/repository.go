package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItem struct {
	Name     string
	Quantity int32
	Price    float64
}

type Order struct {
	ID        int64
	UserID    int64
	Status    string
	Total     float64
	Items     []OrderItem
	CreatedAt time.Time
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, userID int64, items []OrderItem) (*Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var total float64
	for _, it := range items {
		total += float64(it.Quantity) * it.Price
	}

	o := &Order{}
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_id, total) VALUES ($1, $2)
		 RETURNING id, user_id, status, total, created_at`,
		userID, total,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, name, quantity, price) VALUES ($1, $2, $3, $4)`,
			o.ID, it.Name, it.Quantity, it.Price,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	o.Items = items
	return o, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*Order, error) {
	o := &Order{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, status, total, created_at FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT name, quantity, price FROM order_items WHERE order_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.Name, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}
	return o, nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64, page, pageSize int32) ([]*Order, int32, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, status, total, created_at FROM orders
		 WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		o := &Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	var total int32
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID).Scan(&total)
	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, s string) (*Order, error) {
	o := &Order{}
	err := r.db.QueryRow(ctx,
		`UPDATE orders SET status=$1 WHERE id=$2
		 RETURNING id, user_id, status, total, created_at`,
		s, id,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt)
	return o, err
}
