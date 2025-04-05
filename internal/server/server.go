package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/QuizWars-Ecosystem/go-common/pkg/clients"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/handler"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/service"
	"github.com/QuizWars-Ecosystem/users-service/internal/apis/store"

	"github.com/DavidMovas/gopherbox/pkg/closer"
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	"github.com/QuizWars-Ecosystem/go-common/pkg/consul"
	"github.com/QuizWars-Ecosystem/go-common/pkg/log"
	"github.com/QuizWars-Ecosystem/users-service/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var _ abstractions.Server = (*Server)(nil)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	consul     *consul.Consul
	logger     *log.Logger
	cfg        *config.Config
	closer     *closer.Closer
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	cl := closer.NewCloser()

	logger := log.NewLogger(cfg.Local, cfg.LogLevel)
	cl.PushIO(logger)

	consulManager, err := consul.NewConsul(cfg.ConsulURL, cfg.Name, cfg.Address, cfg.GRPCPort, logger)
	if err != nil {
		logger.Zap().Error("error initializing consul manager", zap.Error(err))
		return nil, fmt.Errorf("error initializing consul manager: %w", err)
	}

	cl.Push(consulManager.Stop)

	db, err := clients.NewPostgresClient(ctx, cfg.Postgres.URL, nil)
	if err != nil {
		logger.Zap().Error("error initializing postgres client", zap.Error(err))
		return nil, fmt.Errorf("error initializing postgres client: %w", err)
	}

	storage := store.NewStore(db, logger.Zap())
	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration, cfg.JWT.RefreshExpiration)
	srv := service.NewService(storage, logger.Zap())
	hand := handler.NewHandler(srv, jwtService, logger.Zap())

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())

	healthServer := health.NewServer()
	healthServer.SetServingStatus(fmt.Sprintf("%s-%d", cfg.Name, cfg.GRPCPort), grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	cl.PushNE(healthServer.Shutdown)

	userspb.RegisterUsersAdminServiceServer(grpcServer, hand)
	userspb.RegisterUsersAuthServiceServer(grpcServer, hand)
	userspb.RegisterUsersProfileServiceServer(grpcServer, hand)
	userspb.RegisterUsersSocialServiceServer(grpcServer, hand)

	return &Server{
		grpcServer: grpcServer,
		consul:     consulManager,
		logger:     logger,
		cfg:        cfg,
		closer:     cl,
	}, nil
}

func (s *Server) Start() error {
	z := s.logger.Zap()

	z.Info("Starting server", zap.String("name", s.cfg.Name), zap.Int("port", s.cfg.GRPCPort))

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPCPort))
	if err != nil {
		z.Error("Failed to start listener", zap.String("name", s.cfg.Name), zap.Int("port", s.cfg.GRPCPort), zap.Error(err))
		return err
	}

	s.closer.PushIO(s.listener)

	err = s.consul.RegisterService()
	if err != nil {
		z.Error("Failed to register service in consul registry", zap.String("name", s.cfg.Name), zap.Error(err))
		return err
	}

	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	z := s.logger.Zap()
	z.Info("Shutting down server gracefully", zap.String("name", s.cfg.Name))

	stopChan := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopChan)
	}()

	select {
	case <-stopChan:
	case <-ctx.Done():
		z.Warn("Graceful shutdown timed out, forcing stop")
		s.grpcServer.Stop()
	}

	if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("shutting down listener: %w", err)
	}

	if err := s.logger.Close(); err != nil {
		return fmt.Errorf("error closing logger: %w", err)
	}

	return s.closer.Close(ctx)
}
