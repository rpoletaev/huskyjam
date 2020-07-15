package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Address            string `envconfig:"ADDRESS"`
	MaxIdle            int    `envconfig:"MAX_IDLE"`
	IdleTimeoutSeconds int    `envconfig:"IDLE_TIME"`
}

type Backend struct {
	*Config
	Pool *redis.Pool
}

func (b *Backend) Connect() error {
	b.Pool = &redis.Pool{
		MaxIdle:     b.MaxIdle,
		IdleTimeout: time.Duration(time.Duration(b.IdleTimeoutSeconds) * time.Second),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", b.Address)
		},
	}
	return nil
}

func (b *Backend) Close() error {
	return b.Pool.Close()
}
