package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/grpcx/telemetry"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	manager "github.com/QuizWars-Ecosystem/go-common/pkg/config"
	grpccommon "github.com/QuizWars-Ecosystem/go-common/pkg/grpcx/metrics"
	"github.com/QuizWars-Ecosystem/users-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc/reflection"

	"github.com/QuizWars-Ecosystem/go-common/pkg/clients"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	usersv1 "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
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
	grpcServer   *grpc.Server
	httpServer   *http.Server
	grpcListener net.Listener
	httpListener net.Listener
	consul       *consul.Consul
	logger       *log.Logger
	manager      *manager.Manager[config.Config]
	closer       *closer.Closer
}

func NewServer(ctx context.Context, manager *manager.Manager[config.Config]) (*Server, error) {
	cl := closer.NewCloser()
	cfg := manager.Config()

	logger := log.NewLogger(cfg.Local, cfg.Logger.Level)
	cl.PushIO(logger)

	manager.Subscribe(logger.SectionKey(), func(cfg *config.Config) error { return logger.UpdateConfig(cfg.Logger) })

	consulManager, err := consul.NewConsul(cfg.ConsulURL, cfg.Name, cfg.Address, cfg.GRPCPort, logger)
	if err != nil {
		logger.Zap().Error("error initializing consul manager", zap.Error(err))
		return nil, fmt.Errorf("error initializing consul manager: %w", err)
	}

	cl.Push(consulManager.Stop)

	provider, err := telemetry.NewTracerProvider(ctx, cfg.Name, cfg.Telemetry.URL)
	if err != nil {
		logger.Zap().Error("error initializing telemetry tracer", zap.Error(err))
	}

	cl.PushCtx(provider.Shutdown)

	db, err := clients.NewPostgresClient(ctx, cfg.Postgres.URL,
		clients.NewPostgresOptions(cfg.Postgres.URL).
			WithConnectTimeout(time.Second*10).
			WithTracerProvider(provider),
	)
	if err != nil {
		logger.Zap().Error("error initializing postgres client", zap.Error(err))
		return nil, fmt.Errorf("error initializing postgres client: %w", err)
	}

	grpcprometheus.EnableHandlingTimeHistogram()

	jwtService := jwt.NewService(cfg.JWT)
	manager.Subscribe(jwtService.SectionKey(), func(cfg *config.Config) error { return jwtService.UpdateConfig(cfg.JWT) })

	storage := store.NewStore(db, logger.Zap())
	srv := service.NewService(storage, logger.Zap())
	hand := handler.NewHandler(srv, jwtService, logger.Zap())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcrecovery.UnaryServerInterceptor(),
			grpccommon.ServerMetricsInterceptor(),
			grpcprometheus.UnaryServerInterceptor,
		),
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(provider),
			),
		),
	)

	healthServer := health.NewServer()
	healthServer.SetServingStatus(fmt.Sprintf("%s-%d", cfg.Name, cfg.GRPCPort), grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	cl.PushNE(healthServer.Shutdown)

	usersv1.RegisterUsersAdminServiceServer(grpcServer, hand)
	usersv1.RegisterUsersAuthServiceServer(grpcServer, hand)
	usersv1.RegisterUsersProfileServiceServer(grpcServer, hand)
	usersv1.RegisterUsersSocialServiceServer(grpcServer, hand)

	metrics.Initialize()

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler: metricsMux,
	}

	if cfg.Local {
		reflection.Register(grpcServer)
	}

	return &Server{
		grpcServer: grpcServer,
		httpServer: metricsServer,
		consul:     consulManager,
		logger:     logger,
		manager:    manager,
		closer:     cl,
	}, nil
}

func (s *Server) Start() error {
	z := s.logger.Zap()

	cfg := s.manager.Config()
	s.manager.Watch(z)

	group := errgroup.Group{}

	group.Go(func() error {
		z.Info("Starting metrics server", zap.Int("port", cfg.Metrics.Port))

		var err error
		s.httpListener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Metrics.Port))
		if err != nil {
			z.Error("Error starting metrics server", zap.Error(err))
			return err
		}

		s.closer.PushIO(s.httpListener)

		return s.httpServer.Serve(s.httpListener)
	})

	group.Go(func() error {
		z.Info("Starting server", zap.String("name", cfg.Name), zap.Int("port", cfg.GRPCPort))

		var err error
		s.grpcListener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
		if err != nil {
			z.Error("Failed to start grpcListener", zap.String("name", cfg.Name), zap.Int("port", cfg.GRPCPort), zap.Error(err))
			return err
		}

		s.closer.PushIO(s.grpcListener)

		return s.grpcServer.Serve(s.grpcListener)
	})

	err := s.consul.RegisterService()
	if err != nil {
		z.Error("Failed to register service in consul registry", zap.String("name", cfg.Name), zap.Error(err))
		return err
	}

	return group.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	z := s.logger.Zap()
	cfg := s.manager.Config()

	z.Info("Shutting down server gracefully", zap.String("name", cfg.Name))

	stopChan := make(chan struct{})
	go func() {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			z.Error("Error shutting down metrics server", zap.Error(err))
		}

		s.grpcServer.GracefulStop()
		close(stopChan)
	}()

	select {
	case <-stopChan:
	case <-ctx.Done():
		z.Warn("Graceful shutdown timed out, forcing stop")
		s.grpcServer.Stop()
	}

	if err := s.grpcListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("shutting down grpc listener: %w", err)
	}

	if err := s.httpListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("shutting down http listener: %w", err)
	}

	if err := s.logger.Close(); err != nil {
		return fmt.Errorf("error closing logger: %w", err)
	}

	return s.closer.Close(ctx)
}
