// +build wireinject

package main

import (
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"github.com/rpoletaev/huskyjam/http"
	"github.com/rpoletaev/huskyjam/http/bcrypt"
	"github.com/rpoletaev/huskyjam/internal"
	"github.com/rpoletaev/huskyjam/internal/pg"
	"github.com/rpoletaev/huskyjam/internal/redis"
	"github.com/rpoletaev/huskyjam/pkg/auth"
	"github.com/rpoletaev/huskyjam/pkg/auth/jwt"
	"github.com/rs/zerolog"
)

func providePostgres(c *Config) *pg.Store {
	return &pg.Store{
		Config: c.PG,
	}
}

func provideJWT(c *Config) *jwt.Tokens {
	return jwt.New(c.JWT, nil, nil)
}

func provideRedis(c *Config) *redis.Backend {
	return &redis.Backend{Config: c.Redis}
}

func provideAPIConfig(c *Config) *http.Config {
	return c.HTTP
}

func providePassHelper() *bcrypt.PassManager {
	return &bcrypt.PassManager{}
}

func provideGoodsHandler(logger zerolog.Logger, store internal.Store) *http.GoodsHandler {
	return &http.GoodsHandler{
		Store: store.Goods(),
		Log:   logger.With().Str("service", "goods").Logger(),
	}
}

func provideAccountsHandler(logger zerolog.Logger, store internal.Store, kvstore internal.KVStore, passHelper http.PassHashHelper, tokensHelper auth.Tokens) *http.AccountHandler {
	return &http.AccountHandler{
		Store:          store.Accounts(),
		PassHashHelper: passHelper,
		RefreshRepo:    kvstore.Tokens(),
		Auth:           tokensHelper,
		Log:            logger.With().Str("service", "account").Logger(),
	}
}
func provideApp(logger zerolog.Logger, c *Config) *App {
	wire.Build(
		provideJWT,
		providePostgres,
		provideRedis,
		provideAPIConfig,
		providePassHelper,
		provideGoodsHandler,
		provideAccountsHandler,
		wire.Struct(new(App), "*"),
		wire.Struct(new(http.Api), "*"),
		wire.Bind(new(http.PassHashHelper), new(*bcrypt.PassManager)),
		wire.Bind(new(internal.Store), new(*pg.Store)),
		wire.Bind(new(internal.KVStore), new(*redis.Backend)),
		wire.Bind(new(auth.Tokens), new(*jwt.Tokens)),
	)

	return &App{}
}
