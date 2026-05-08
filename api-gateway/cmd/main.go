package main

import (
	"log"
	"net/http"
	"os"

	bookingpb "github.com/final/gen/booking"
	orderpb "github.com/final/gen/order"
	userpb "github.com/final/gen/user"
	"github.com/final/api-gateway/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	userConn, err := grpc.NewClient(getenv("USER_SERVICE_ADDR", "localhost:50051"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("user-service connect: %v", err)
	}
	defer userConn.Close()

	orderConn, err := grpc.NewClient(getenv("ORDER_SERVICE_ADDR", "localhost:50052"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("order-service connect: %v", err)
	}
	defer orderConn.Close()

	bookingConn, err := grpc.NewClient(getenv("BOOKING_SERVICE_ADDR", "localhost:50053"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("booking-service connect: %v", err)
	}
	defer bookingConn.Close()

	h := handler.New(
		userpb.NewUserServiceClient(userConn),
		orderpb.NewOrderServiceClient(orderConn),
		bookingpb.NewBookingServiceClient(bookingConn),
	)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := getenv("HTTP_ADDR", ":8080")
	log.Printf("api-gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
