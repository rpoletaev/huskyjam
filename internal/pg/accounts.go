package pg

import (
	"github.com/rpoletaev/huskyjam/internal"
)

const (
	initAccountsSchema = `CREATE TABLE IF NOT EXISTS accounts (
		id integer PRIMARY KEY,
		email text NOT NULL,
		created_at timestamp DEFAULT current_timestamp
		deleted_at timestamp 
		CONSTRAINT unique_acc_email UNIQUE(email),
		)`
)

type AccountsRepository Store

func (s *Store) Accounts() internal.AccountsRepository {
	return (*AccountsRepository)(s)
}

var _ internal.AccountsRepository = (*AccountsRepository)(nil)

func (s *AccountsRepository) Init() {
	s.db.MustExec(initAccountsSchema)
}

func (s *AccountsRepository) Create(acc *internal.Account) error {

	if _, err := s.db.Exec("INSERT INTO accounts (email, pass) VALUES (:email, :pass)", acc); err != nil {
		if uniqueViolation(err) {
			return internal.ErrAlreadyExists
		}

		return accessDB(err)
	}

	return nil
}

func (s *AccountsRepository) GetByEmail(email string) (*internal.Account, error) {
	acc := &internal.Account{}
	if err := s.db.Get(acc, "SELECT * FROM accounts WHERE email = $1 and deleted_at IS NOT NULL", email); err != nil {
		if notFound(err) {
			return nil, internal.ErrNotFound
		}

		return nil, accessDB(err)
	}

	return acc, nil
}
