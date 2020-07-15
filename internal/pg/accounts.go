package pg

import (
	"github.com/rpoletaev/huskyjam/internal"
)

const (
	initAccountsSchema = `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		email text NOT NULL,
		pass text NOT NULL,
		created_at timestamp DEFAULT current_timestamp,
		deleted_at timestamp, 
		CONSTRAINT unique_acc_email UNIQUE(email)
		)`
)

type AccountsRepository Store

func (s *Store) Accounts() internal.AccountsRepository {
	return (*AccountsRepository)(s)
}

var _ internal.AccountsRepository = (*AccountsRepository)(nil)

func (s *AccountsRepository) Init() error {
	if _, err := s.db.Exec(initAccountsSchema); err != nil {
		return err
	}

	return nil
}

func (s *AccountsRepository) Create(acc *internal.Account) error {

	if _, err := s.db.NamedExec("INSERT INTO accounts (email, pass) VALUES (:email, :pass)", acc); err != nil {
		if uniqueViolation(err) {
			return internal.ErrAlreadyExists
		}

		return accessDB(err)
	}

	return nil
}

func (s *AccountsRepository) GetByEmail(email string) (*internal.Account, error) {
	acc := &internal.Account{}
	if err := s.db.Get(acc, "SELECT * FROM accounts WHERE email = $1 and deleted_at IS NULL", email); err != nil {
		if notFound(err) {
			return nil, internal.ErrNotFound
		}

		return nil, accessDB(err)
	}

	return acc, nil
}
