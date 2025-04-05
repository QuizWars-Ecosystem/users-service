package config

import (
	"time"

	def "github.com/QuizWars-Ecosystem/go-common/pkg/config"
	"github.com/QuizWars-Ecosystem/users-service/internal/config"
	"github.com/testcontainers/testcontainers-go"
)

type TestConfig struct {
	ServiceConfig *config.Config
	Postgres      *PostgresConfig
	Network       *testcontainers.DockerNetwork
}

func NewTestConfig() *TestConfig {
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
			Postgres: config.PostgresConfig{
				URL: "postgres:5432",
			},
			Redis: config.RedisConfig{
				URL: "redis:6379",
			},
		},
		Postgres: &PostgresConfig{
			Name:     "postgres",
			Image:    "postgres:17.2-alpine",
			Username: "user",
			Password: "pass",
			DBName:   "users",
			Host:     "localhost",
			Port:     5432,
		},
	}
}

type PostgresConfig struct {
	Name     string
	Image    string
	Username string
	Password string
	DBName   string
	Host     string
	Port     int
}
