package auth

import (
	"context"
)

// Tokens provides auth tokens
type Tokens interface {
	Verify(token string) (*SystemClaims, error)
	SignToken(src *SystemClaims) (string, error)
}

type SystemClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (claims *SystemClaims) GetID() uint {
	if claims == nil {
		return 0
	}
	return claims.ID
}

type claimsKey struct{}

// WithClaims injects claims into context
func WithClaims(ctx context.Context, claims *SystemClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// Get extracts claims from context
func Get(ctx context.Context) *SystemClaims {
	res, _ := ctx.Value(claimsKey{}).(*SystemClaims)
	return res
}
