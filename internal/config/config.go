package config

import "github.com/QuizWars-Ecosystem/go-common/pkg/config"

type Config struct {
	config.DefaultServiceConfig
	PostgresConfig `envPrefix:"POSTGRES_"`
	RedisConfig    `envPrefix:"REDIS_"`
}

type PostgresConfig struct {
	URL string `env:"URL"`
}

type RedisConfig struct {
	URL string `env:"URL"`
}
