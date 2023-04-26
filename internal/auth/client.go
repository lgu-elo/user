package auth

import (
	"github.com/lgu-elo/auth/pkg/pb"
	"github.com/lgu-elo/user/internal/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client pb.AuthClient

func NewClient(cfg *config.Cfg) (Client, error) {
	conn, err := grpc.Dial(cfg.Services.Auth.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "canot")
	}

	client := pb.NewAuthClient(conn)

	return client, nil
}
