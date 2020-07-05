package pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/rpoletaev/huskyjam/internal"
)

// Config for postgres connection
type Config struct {
	Driver string `envconfig:"DRIVER"`
	URI    string `envconfig:"URI"`
}

// Store implements internal Store with postgres db
type Store struct {
	*Config
	db *sqlx.DB
}

var _ internal.Store = (*Store)(nil)

// Connect fulfill db connection
func (s *Store) Connect() error {
	db, err := sqlx.Connect(s.Driver, s.URI)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
