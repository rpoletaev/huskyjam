package jwt

import (
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/rpoletaev/huskyjam/pkg/auth"
)

// Config struct to configure jwt
type Config struct {
	Secret           string `envconfig:"SECRET" default:"test secret"`
	AccessTTLMinutes int    `envconfig:"ACCESS_TTL_MINUTES"`
}

type UUIDGenerator interface {
	Generate() string
}

type CurrentTimeGenerator interface {
	Now() time.Time
}

// Tokens implements basic services functions and middlewares
type Tokens struct {
	SigningMethod jwt.SigningMethod
	UUID          UUIDGenerator `wire:"-"`
	TimeToLive    time.Duration
	Secret        string
	TimeGetter    CurrentTimeGenerator `wire:"-"`
	Debug         bool

	mu sync.RWMutex
}

var _ auth.Tokens = (*Tokens)(nil)

// Claims struct for jwt claims
type Claims struct {
	jwt.StandardClaims
	*auth.SystemClaims
}

type xidUUIDGenerator struct{}

func (g xidUUIDGenerator) Generate() string {
	return xid.New().String()
}

type defaultTimeGetter struct{}

func (g defaultTimeGetter) Now() time.Time {
	return time.Now()
}

// New creates new concrete services.Tokens instance
func New(c *Config, uuidGen UUIDGenerator, timeGetter CurrentTimeGenerator) *Tokens {
	defineUUIDGenerator := func() UUIDGenerator {
		if uuidGen == nil {
			return xidUUIDGenerator{}
		}

		return uuidGen
	}

	defineTimeGetter := func() CurrentTimeGenerator {
		if timeGetter == nil {
			return defaultTimeGetter{}
		}

		return timeGetter
	}

	return &Tokens{
		SigningMethod: jwt.SigningMethodHS256,
		UUID:          defineUUIDGenerator(),
		TimeGetter:    defineTimeGetter(),
		Secret:        c.Secret,
		TimeToLive:    time.Duration(c.AccessTTLMinutes) * time.Minute,
	}
}

// Verify token
func (p *Tokens) Verify(token string) (*auth.SystemClaims, error) {
	token = strings.TrimSpace(token)
	var res Claims
	parser := &jwt.Parser{SkipClaimsValidation: true}
	_, err := parser.ParseWithClaims(token, &res, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.Secret), nil
	})
	if err != nil {
		return nil, auth.ErrUnauthorized
	}
	switch {
	case !res.VerifyExpiresAt(jwt.TimeFunc().Unix(), true):
		return nil, auth.ErrTokenExpired
	default:
		if err = res.Valid(); err != nil {
			return nil, auth.ErrInvalidToken
		}
	}
	return res.SystemClaims, nil
}

// SignToken creates string from user
func (p *Tokens) SignToken(u *auth.SystemClaims) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        p.UUID.Generate(),
			ExpiresAt: p.TimeGetter.Now().Add(p.TimeToLive).Unix(),
		},
		SystemClaims: u,
	}
	token := jwt.NewWithClaims(p.SigningMethod, claims)
	res, err := token.SignedString([]byte(p.Secret))
	if err != nil {
		return "", errors.Wrap(err, "sign token")
	}
	return res, nil
}
