package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Backend struct {
	Pool               *redis.Pool
	Address            string
	MaxIdle            int
	IdleTimeoutSeconds int
}

func (b *Backend) Connect() {
	b.Pool = &redis.Pool{
		MaxIdle:     b.MaxIdle,
		IdleTimeout: time.Duration(time.Duration(b.IdleTimeoutSeconds) * time.Second),
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", b.Address) },
	}
}

func (b *Backend) Close() error {
	return b.Pool.Close()
}
