package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/DavidMovas/gopherbox/pkg/closer"
	"github.com/QuizWars-Ecosystem/go-common/pkg/clients"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/log"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/handler"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/service"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/store"
	"github.com/QuizWars-Ecosystem/users-service/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TestServer struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     *log.Logger
	cfg        *config.Config
	closer     *closer.Closer
}

func NewTestServer(ctx context.Context, cfg *config.Config) (*TestServer, error) {
	cl := closer.NewCloser()

	logger := log.NewLogger(cfg.Local, cfg.LogLevel)
	cl.PushIO(logger)

	dbOpts := clients.NewPostgresOptions(cfg.Postgres.URL)
	dbOpts.WithConnectTimeout(time.Second * 20)

	db, err := clients.NewPostgresClient(ctx, cfg.Postgres.URL, dbOpts)
	if err != nil {
		logger.Zap().Error("error initializing postgres client", zap.Error(err))
		return nil, fmt.Errorf("error initializing postgres client: %w", err)
	}

	storage := store.NewStore(db, logger.Zap())
	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration, cfg.JWT.RefreshExpiration)
	srv := service.NewService(storage, logger.Zap())
	hand := handler.NewHandler(srv, jwtService, logger.Zap())

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())

	userspb.RegisterUsersAdminServiceServer(grpcServer, hand)
	userspb.RegisterUsersAuthServiceServer(grpcServer, hand)
	userspb.RegisterUsersProfileServiceServer(grpcServer, hand)
	userspb.RegisterUsersSocialServiceServer(grpcServer, hand)

	return &TestServer{
		grpcServer: grpcServer,
		logger:     logger,
		cfg:        cfg,
		closer:     cl,
	}, nil
}

func (s *TestServer) Start() error {
	z := s.logger.Zap()

	z.Info("Starting server", zap.String("name", s.cfg.Name), zap.Int("port", s.cfg.GRPCPort))

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPCPort))
	if err != nil {
		z.Error("Failed to start listener", zap.String("name", s.cfg.Name), zap.Int("port", s.cfg.GRPCPort), zap.Error(err))
		return err
	}

	return s.grpcServer.Serve(s.listener)
}

func (s *TestServer) Shutdown(ctx context.Context) error {
	z := s.logger.Zap()

	z.Info("Shutting down server", zap.String("name", s.cfg.Name))

	s.grpcServer.GracefulStop()

	err := s.listener.Close()
	if err != nil {
		z.Error("Error shutting down listener", zap.String("name", s.cfg.Name), zap.Error(err))
	}

	return s.closer.Close(ctx)
}
