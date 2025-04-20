package config

import (
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/log"

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
			ServiceConfig: &def.ServiceConfig{
				Name:         "users-service",
				Address:      "users_address",
				Local:        true,
				GRPCPort:     50051,
				StartTimeout: time.Second * 30,
				StopTimeout:  time.Second * 30,
				ConsulURL:    "consul:8500",
			},
			Logger: &log.Config{
				Level: "debug",
			},
			JWT: &jwt.Config{
				Secret:            "secret",
				AccessExpiration:  time.Hour,
				RefreshExpiration: time.Hour,
			},
		},
		Postgres: &postgresCfg,
	}
}
