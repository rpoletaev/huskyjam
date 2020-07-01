package redis

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
	"github.com/rpoletaev/huskyjam/pkg/auth"
	"github.com/rs/xid"
)

const tokenKeyTemplate = "token:%s"

// TokensRepository implements internal.TokenRepository
type TokensRepository Backend

// Tokens returns internal.TokensRepository
func (repo *Backend) Tokens() internal.TokensRepository {
	return (*TokensRepository)(repo)
}

func (t *TokensRepository) uuid() string {
	return xid.New().String()
}

var _ internal.TokensRepository = (*TokensRepository)(nil)

// New creates refresh token info
func (t *TokensRepository) New(val *auth.SystemClaims) (string, error) {
	c := t.Pool.Get()
	defer c.Close()

	bts, err := json.Marshal(val)
	if err != nil {
		return "", errors.Wrap(err, "marshal claims")
	}

	token := t.uuid()
	_, err = c.Do("set", fmt.Sprintf(tokenKeyTemplate, token), bts)

	return token, errors.Wrap(err, "save token")
}

// Get returns claims by refresh token
func (t *TokensRepository) Get(token string) (*auth.SystemClaims, error) {
	c := t.Pool.Get()
	defer c.Close()

	res, err := redis.String(c.Do("get", fmt.Sprintf(tokenKeyTemplate, token)))
	if err != nil {
		if err == redis.ErrNil {
			err = internal.ErrNotFound
		}

		return nil, errors.Wrap(err, "get token")
	}

	claims := &auth.SystemClaims{}
	return claims, errors.Wrap(json.Unmarshal([]byte(res), claims), "unmarshal token")
}

// Delete deletes refresh token info
func (t *TokensRepository) Delete(token string) error {
	c := t.Pool.Get()
	defer c.Close()

	_, err := c.Do("del", fmt.Sprintf(tokenKeyTemplate, token))
	return errors.Wrap(err, "delete token")
}
