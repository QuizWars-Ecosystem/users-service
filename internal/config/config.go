package config

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/config"
	"github.com/QuizWars-Ecosystem/go-common/pkg/jwt"
	"github.com/QuizWars-Ecosystem/go-common/pkg/log"
)

type Config struct {
	*config.ServiceConfig `mapstructure:"service"`
	Logger                *log.Config     `mapstructure:"logger"`
	JWT                   *jwt.Config     `mapstructure:"jwt"`
	Postgres              *PostgresConfig `mapstructure:"postgres"`
	Redis                 *RedisConfig    `mapstructure:"redis"`
}

type PostgresConfig struct {
	URL string `mapstructure:"url"`
}

type RedisConfig struct {
	URL string `mapstructure:"url"`
}
