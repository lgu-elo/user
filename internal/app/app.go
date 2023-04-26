package app

import (
	"fmt"
	"net"
	"sync"

	"github.com/lgu-elo/gateway/pkg/logger"
	"github.com/lgu-elo/gateway/pkg/rpc"
	"github.com/lgu-elo/user/internal/adapters/database"
	"github.com/lgu-elo/user/internal/config"
	"github.com/lgu-elo/user/internal/server"
	"github.com/lgu-elo/user/internal/user"
	"github.com/lgu-elo/user/pkg/pb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	fxlogrus "github.com/takt-corp/fx-logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.WithLogger(func(log *logrus.Logger) fxevent.Logger {
			return &fxlogrus.LogrusLogger{Logger: log}
		}),
		fx.Provide(
			logger.New,
			config.New,
			database.New,
			func() *sync.Mutex {
				var mu sync.Mutex
				return &mu
			},

			fx.Annotate(user.NewStorage, fx.As(new(user.Repository))),
			fx.Annotate(user.NewService, fx.As(new(user.IService))),
			fx.Annotate(user.NewHandler, fx.As(new(user.IHandler))),
			//fx.Annotate(auth.NewClient, fx.As(new(auth.Client))),

			server.NewAPI,

			func(logger *logrus.Logger) *grpc.Server {
				return rpc.New(
					rpc.WithLoggerLogrus(logger),
				)
			},
		),
		fx.Invoke(serve),
	)
}

func serve(logger *logrus.Logger, srv *grpc.Server, cfg *config.Cfg, api *server.API) {
	pb.RegisterUserServiceServer(srv, api)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		logger.Error(errors.Wrap(err, "cannot bind server"))
		return
	}

	if errServe := srv.Serve(lis); errServe != nil {
		logger.Error(errors.Wrap(err, "cannot lsiten server"))
		return
	}
}
