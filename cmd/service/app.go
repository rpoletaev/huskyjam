package main

import (
	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/http"
	"github.com/rpoletaev/huskyjam/internal/pg"
	"github.com/rpoletaev/huskyjam/internal/redis"
	"github.com/rpoletaev/huskyjam/pkg/auth/jwt"
	"github.com/rs/zerolog"
)

// App structs combines concrete side systems and handles their connections
type App struct {
	Store   *pg.Store
	KVStore *redis.Backend
	Tokens  *jwt.Tokens
	API     *http.Api
	Log     zerolog.Logger
}

// Connect opens connections to side systems
func (app *App) Connect() error {
	if err := app.Store.Connect(); err != nil {
		return errors.Wrap(err, "connect to postgress")
	}

	app.Log.Info().Msg("connected to db")

	if err := app.Store.Accounts().Init(); err != nil {
		return errors.Wrap(err, "run migrations for accounts")
	}

	app.Log.Info().Msg("accounts migrations executed")

	if err := app.Store.Goods().Init(); err != nil {
		return errors.Wrap(err, "run migrations for goods")
	}

	app.Log.Info().Msg("goods migrations executed")

	if err := app.KVStore.Connect(); err != nil {
		return errors.Wrap(err, "connect to redis")
	}

	app.Log.Info().Msg("connected to redis")

	return nil
}

// Close closes connections to side systems
func (app *App) Close() error {
	if err := app.Store.Close(); err != nil {
		return errors.Wrap(err, "close postgress connection")
	}

	if err := app.KVStore.Close(); err != nil {
		return errors.Wrap(err, "close redis connection")
	}
	return nil
}
