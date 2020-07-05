package main

import (
	"github.com/rpoletaev/huskyjam/http"
	"github.com/rpoletaev/huskyjam/internal/pg"
	"github.com/rpoletaev/huskyjam/internal/redis"
	"github.com/rpoletaev/huskyjam/pkg/auth/jwt"

	_ "github.com/joho/godotenv/autoload" // preload .env
	"github.com/kelseyhightower/envconfig"

	"github.com/pkg/errors"
)

// Config for application
type Config struct {
	Debug bool          `envconfig:"DEBUG" default:"true"`
	PG    *pg.Config    `envconfig:"DB"`
	Redis *redis.Config `envconfig:"REDIS"`
	JWT   *jwt.Config   `envconfig:"JWT"`
	HTTP  *http.Config  `envconfig:"HTTP"`
}

// LoadConfig parses env and loads config
func LoadConfig(prefix string) (*Config, error) {
	res := &Config{}
	return res, errors.Wrap(envconfig.Process(prefix, res), "parse config")
}
