package server

import (
	"context"
	"encoding/json"

	bookingpb "github.com/final/gen/booking"
	"github.com/final/booking-service/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingServer struct {
	bookingpb.UnimplementedBookingServiceServer
	repo    *repository.BookingRepository
	channel *amqp.Channel
}

func New(repo *repository.BookingRepository, ch *amqp.Channel) *BookingServer {
	return &BookingServer{repo: repo, channel: ch}
}

func (s *BookingServer) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.BookingResponse, error) {
	b, err := s.repo.Create(ctx, req.UserId, req.Resource, req.StartTime, req.EndTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create booking: %v", err)
	}
	resp := toProto(b)
	s.publish("booking.created", resp)
	return resp, nil
}

func (s *BookingServer) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.BookingResponse, error) {
	b, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "booking not found: %v", err)
	}
	return toProto(b), nil
}

func (s *BookingServer) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	bookings, total, err := s.repo.ListByUser(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list bookings: %v", err)
	}
	var items []*bookingpb.BookingResponse
	for _, b := range bookings {
		items = append(items, toProto(b))
	}
	return &bookingpb.ListBookingsResponse{Bookings: items, Total: total}, nil
}

func (s *BookingServer) UpdateBookingStatus(ctx context.Context, req *bookingpb.UpdateBookingStatusRequest) (*bookingpb.BookingResponse, error) {
	b, err := s.repo.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update booking status: %v", err)
	}
	resp := toProto(b)
	s.publish("booking.status_updated", resp)
	return resp, nil
}

func (s *BookingServer) publish(routingKey string, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	s.channel.Publish("bookings", routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func toProto(b *repository.Booking) *bookingpb.BookingResponse {
	return &bookingpb.BookingResponse{
		Id:        b.ID,
		UserId:    b.UserID,
		Resource:  b.Resource,
		StartTime: b.StartTime.String(),
		EndTime:   b.EndTime.String(),
		Status:    b.Status,
		CreatedAt: b.CreatedAt.String(),
	}
}
