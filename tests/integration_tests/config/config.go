package config

import (
	"time"

	test "github.com/QuizWars-Ecosystem/go-common/pkg/testing/config"

	def "github.com/QuizWars-Ecosystem/go-common/pkg/config"
	"github.com/QuizWars-Ecosystem/users-service/internal/config"
)

type TestConfig struct {
	ServiceConfig *config.Config
	Postgres      *test.PostgresConfig
}

func NewTestConfig() *TestConfig {
	postgresCfg := test.DefaultPostgresConfig()

	return &TestConfig{
		ServiceConfig: &config.Config{
			DefaultServiceConfig: def.DefaultServiceConfig{
				Name:            "users-service",
				Address:         "users_address",
				Local:           true,
				LogLevel:        "debug",
				GRPCPort:        50051,
				StartTimeout:    time.Second * 30,
				ShutdownTimeout: time.Second * 30,
				ConsulURL:       "consul:8500",
			},
			JWT: config.JWTConfig{
				Secret:            "secret",
				AccessExpiration:  time.Hour,
				RefreshExpiration: time.Hour,
			},
		},
		Postgres: &postgresCfg,
	}
}
