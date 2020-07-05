package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rpoletaev/huskyjam/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	baseLogger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := baseLogger.With().
		Dict("meta", zerolog.Dict().
			Str("version", pkg.Version).
			Str("git_sha", pkg.GitSHA).
			Str("timestamp", pkg.Timestamp),
		).Logger()

	cfg, err := LoadConfig("")
	if err != nil {
		logger.Fatal().Err(err).Msg("on read config")
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	app := provideApp(logger, cfg)
	if err := app.Connect(); err != nil {
		logger.Fatal().Err(err).Msg("on app connect")
	}

	defer app.Close()

	go func() {
		if err := app.API.Server().ListenAndServe(); err != nil || err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("on api handle")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("shutdown server ...")
	if err := app.API.Server().Shutdown(context.Background()); err != nil {
		logger.Error().Err(err).Msg("on shutdown server")
	}

}
