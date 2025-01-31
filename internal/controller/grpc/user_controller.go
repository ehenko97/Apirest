package grpc

import (
	"context"
	"fmt"
	"github.com/ehenko97/apirest/internal/entity"
	service "github.com/ehenko97/apirest/internal/service"
	pb "github.com/ehenko97/apirest/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserController struct {
	service service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (h *UserController) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := &entity.User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := h.service.Create(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        int64(createdUser.ID),
			Name:      createdUser.Name,
			Email:     createdUser.Email,
			CreatedAt: timestamppb.New(createdUser.CreatedAt),
			UpdatedAt: timestamppb.New(createdUser.UpdatedAt),
		},
	}, nil
}

func (h *UserController) FindByID(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.service.FindByID(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        int64(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (h *UserController) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user := &entity.User{
		ID:    int(req.Id),
		Name:  req.Name,
		Email: req.Email,
	}

	updatedUser, err := h.service.Update(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        int64(updatedUser.ID),
			Name:      updatedUser.Name,
			Email:     updatedUser.Email,
			CreatedAt: timestamppb.New(updatedUser.CreatedAt),
			UpdatedAt: timestamppb.New(updatedUser.UpdatedAt),
		},
	}, nil
}

func (h *UserController) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.service.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &pb.DeleteUserResponse{}, nil
}

func (h *UserController) FindAll(ctx context.Context, _ *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := h.service.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	response := &pb.ListUsersResponse{}
	for _, user := range users {
		response.Users = append(response.Users, &pb.User{
			Id:        int64(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}
	return response, nil
}
