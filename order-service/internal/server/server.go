package server

import (
	"context"
	"encoding/json"

	orderpb "github.com/final/gen/order"
	"github.com/final/order-service/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	repo    *repository.OrderRepository
	channel *amqp.Channel
}

func New(repo *repository.OrderRepository, ch *amqp.Channel) *OrderServer {
	return &OrderServer{repo: repo, channel: ch}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	var items []repository.OrderItem
	for _, it := range req.Items {
		items = append(items, repository.OrderItem{Name: it.Name, Quantity: it.Quantity, Price: it.Price})
	}
	o, err := s.repo.Create(ctx, req.UserId, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create order: %v", err)
	}
	resp := toProto(o)
	s.publish("order.created", resp)
	return resp, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	o, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return toProto(o), nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	orders, total, err := s.repo.ListByUser(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list orders: %v", err)
	}
	var items []*orderpb.OrderResponse
	for _, o := range orders {
		items = append(items, toProto(o))
	}
	return &orderpb.ListOrdersResponse{Orders: items, Total: total}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	o, err := s.repo.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update order status: %v", err)
	}
	resp := toProto(o)
	s.publish("order.status_updated", resp)
	return resp, nil
}

func (s *OrderServer) publish(routingKey string, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	s.channel.Publish("bookings", routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func toProto(o *repository.Order) *orderpb.OrderResponse {
	resp := &orderpb.OrderResponse{
		Id:        o.ID,
		UserId:    o.UserID,
		Status:    o.Status,
		Total:     o.Total,
		CreatedAt: o.CreatedAt.String(),
	}
	for _, it := range o.Items {
		resp.Items = append(resp.Items, &orderpb.OrderItem{
			Name:     it.Name,
			Quantity: it.Quantity,
			Price:    it.Price,
		})
	}
	return resp
}
