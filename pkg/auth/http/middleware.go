package http

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/pkg/auth"
)

// Основные параметры
var (
	TokenHeader = "X-Auth-Key"
)

type verifier interface {
	Verify(token string) (*auth.SystemClaims, error)
}

func belongsTo(s string, ss []string) bool {
	for _, str := range ss {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}

// WithAuth returns middleware which injects confirmed token into context.
// If there is no token, the context is not populated. Otherwise,
// the token is checked and status 419 is returned if the token is out of date.
func WithAuth(svc verifier, exceptions ...string) func(http.Handler) http.Handler {
	return func(src http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if belongsTo(r.URL.Path, exceptions) {
					src.ServeHTTP(w, r)
					return
				}
				tok := r.Header.Get(TokenHeader)
				claims, err := svc.Verify(tok)
				if errors.Cause(err) == auth.ErrTokenExpired {
					http.Error(w, "token expired", 419)
					return
				} else if err != nil {
					http.Error(w, err.Error(), 401)
					return
				}

				ctx := auth.WithClaims(r.Context(), claims)
				src.ServeHTTP(w, r.WithContext(ctx))
			})
	}
}
