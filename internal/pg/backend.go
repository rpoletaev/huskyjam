package pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
)

// Config for postgres connection
type Config struct {
	Driver string `envconfig:"DRIVER"`
	URI    string `envconfig:"URI"`
}

// Store implements internal Store with postgres db
type Store struct {
	db *sqlx.DB
}

var _ internal.Store = (*Store)(nil)

// Connect fulfill db connection
func (s *Store) Connect() error {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		return errors.Wrap(err, "on connect to postgres")
	}

	s.db = db

	return nil
}
