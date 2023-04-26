package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/lgu-elo/auth/pkg/pb"
	"github.com/lgu-elo/user/internal/auth"
	"github.com/lgu-elo/user/internal/user/model"
)

type (
	IService interface {
		GetAllUsers() ([]*model.User, error)
		GetUserById(id int) (*model.User, error)
		UpdateUser(user *model.User) (*model.User, error)
		DeleteUser(id int) error
		CreateUser(creds *pb.BasicCredentials) error
	}

	service struct {
		repo       Repository
		authClient auth.Client
	}
)

func NewService(repo Repository, authClient auth.Client) IService {
	return &service{repo, authClient}
}

func (s *service) DeleteUser(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllUsers() ([]*model.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "can't get all users")
	}

	return users, nil

}

func (s *service) UpdateUser(user *model.User) (*model.User, error) {
	if err := s.repo.Update(user); err != nil {
		return nil, errors.Wrap(err, "can't update user")
	}

	user, err := s.GetUserById(user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) CreateUser(creds *pb.BasicCredentials) error {
	_, errRegister := s.authClient.Register(context.Background(), creds)
	if errRegister != nil {
		return errors.Wrap(errRegister, "can't register user")
	}

	return nil
}

func (s *service) GetUserById(id int) (*model.User, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	return user, nil
}
