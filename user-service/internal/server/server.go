package server

import (
	"context"

	"github.com/final/gen/user"
	"github.com/final/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func New(repo *repository.UserRepository) *UserServer {
	return &UserServer{repo: repo}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash password: %v", err)
	}
	u, err := s.repo.Create(ctx, req.Name, req.Email, string(hashed))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}
	return toProto(u), nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	u, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	return toProto(u), nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	u, err := s.repo.Update(ctx, req.Id, req.Name, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user: %v", err)
	}
	return toProto(u), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user: %v", err)
	}
	return &userpb.DeleteUserResponse{Success: true}, nil
}

func toProto(u *repository.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.String(),
	}
}
