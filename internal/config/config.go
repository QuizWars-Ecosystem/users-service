package config

import (
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/config"
)

type Config struct {
	config.DefaultServiceConfig
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
}

type JWTConfig struct {
	Secret            string        `env:"SECRET"`
	AccessExpiration  time.Duration `env:"ACCESS_EXPIRATION"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION"`
}

type PostgresConfig struct {
	URL string `env:"URL"`
}

type RedisConfig struct {
	URL string `env:"URL"`
}
