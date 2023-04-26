package user

import (
	"context"

	auth_pb "github.com/lgu-elo/auth/pkg/pb"
	"github.com/lgu-elo/user/internal/user/model"
	"github.com/lgu-elo/user/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type (
	userHandler struct {
		service IService
		log     *logrus.Logger
		server  *grpc.Server
	}

	IHandler interface {
		GetAllUsers(c context.Context, _ *pb.Empty) (*pb.UsersList, error)
		GetUserById(c context.Context, user *pb.UserWithID) (*pb.User, error)
		UpdateUser(c context.Context, user *pb.User) (*pb.User, error)
		DeleteUser(c context.Context, user *pb.UserWithID) (*pb.Empty, error)
		CreateUser(c context.Context, user *pb.User) (*pb.Empty, error)
	}
)

func NewHandler(service IService, log *logrus.Logger, server *grpc.Server) IHandler {
	return &userHandler{service, log, server}
}

func (h *userHandler) GetAllUsers(c context.Context, _ *pb.Empty) (*pb.UsersList, error) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var pbUsers pb.UsersList
	for _, user := range users {
		pbUsers.Users = append(pbUsers.Users, &pb.User{
			Id:    int64(user.ID),
			Login: user.Login,
			Role:  user.Role,
			Name:  user.Name,
		})
	}

	return &pbUsers, nil
}

func (h *userHandler) GetUserById(c context.Context, request *pb.UserWithID) (*pb.User, error) {
	user, err := h.service.GetUserById(int(request.Id))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    int64(user.ID),
		Login: user.Login,
		Role:  user.Role,
		Name:  user.Name,
	}, nil
}

func (h *userHandler) UpdateUser(c context.Context, request *pb.User) (*pb.User, error) {
	user, err := h.service.UpdateUser(&model.User{
		ID:    int(request.Id),
		Login: request.Login,
		Name:  request.Name,
		Role:  request.Role,
	})
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    int64(user.ID),
		Login: user.Login,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (h *userHandler) DeleteUser(c context.Context, user *pb.UserWithID) (*pb.Empty, error) {
	if err := h.service.DeleteUser(int(user.Id)); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (h *userHandler) CreateUser(c context.Context, user *pb.User) (*pb.Empty, error) {
	err := h.service.CreateUser(&auth_pb.BasicCredentials{
		Username: user.Login,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
