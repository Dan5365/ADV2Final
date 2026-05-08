package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Booking struct {
	ID        int64
	UserID    int64
	Resource  string
	StartTime time.Time
	EndTime   time.Time
	Status    string
	CreatedAt time.Time
}

type BookingRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, userID int64, resource, startTime, endTime string) (*Booking, error) {
	b := &Booking{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO bookings (user_id, resource, start_time, end_time)
		 VALUES ($1, $2, $3::timestamptz, $4::timestamptz)
		 RETURNING id, user_id, resource, start_time, end_time, status, created_at`,
		userID, resource, startTime, endTime,
	).Scan(&b.ID, &b.UserID, &b.Resource, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt)
	return b, err
}

func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*Booking, error) {
	b := &Booking{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, resource, start_time, end_time, status, created_at FROM bookings WHERE id = $1`, id,
	).Scan(&b.ID, &b.UserID, &b.Resource, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt)
	return b, err
}

func (r *BookingRepository) ListByUser(ctx context.Context, userID int64, page, pageSize int32) ([]*Booking, int32, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, resource, start_time, end_time, status, created_at
		 FROM bookings WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []*Booking
	for rows.Next() {
		b := &Booking{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.Resource, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt); err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}

	var total int32
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM bookings WHERE user_id = $1`, userID).Scan(&total)
	return bookings, total, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id int64, status string) (*Booking, error) {
	b := &Booking{}
	err := r.db.QueryRow(ctx,
		`UPDATE bookings SET status=$1 WHERE id=$2
		 RETURNING id, user_id, resource, start_time, end_time, status, created_at`,
		status, id,
	).Scan(&b.ID, &b.UserID, &b.Resource, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt)
	return b, err
}
