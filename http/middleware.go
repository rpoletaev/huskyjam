package http

import (
	"net/http"
	"os"
	"time"

	recover "github.com/dre1080/recovr"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	recovery := recover.New()

	logEvent := func(logger *zerolog.Logger, status int) *zerolog.Event {
		if status >= 200 && status < 300 {
			return logger.Info()
		}

		if status >= 400 {
			return logger.Error()
		}

		return logger.Debug()
	}

	chain := alice.New()
	chain = chain.Append(recovery)
	chain = chain.Append(cors.AllowAll().Handler)
	chain = chain.Append(hlog.NewHandler(zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()))
	chain = chain.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		logger := hlog.FromRequest(r)
		logEvent(logger, status).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	chain = chain.Append(hlog.RemoteAddrHandler("ip"))
	chain = chain.Append(hlog.UserAgentHandler("user_agent"))
	chain = chain.Append(hlog.RefererHandler("referer"))
	chain = chain.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return chain.Then(handler)
}
