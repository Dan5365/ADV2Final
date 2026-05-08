package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"

	orderpb "github.com/final/gen/order"
	"github.com/final/order-service/internal/repository"
	"github.com/final/order-service/internal/server"
	migrations "github.com/final/order-service/migrations"
	"google.golang.org/grpc"
)

func main() {
	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/orders?sslmode=disable")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect: %v", err)
	}
	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("goose up: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()

	rabbitURL := getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbitmq dial: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("rabbitmq channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("bookings", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("exchange declare: %v", err)
	}

	addr := getenv("GRPC_ADDR", ":50052")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	repo := repository.New(pool)
	orderpb.RegisterOrderServiceServer(srv, server.New(repo, ch))

	log.Printf("order-service listening on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
