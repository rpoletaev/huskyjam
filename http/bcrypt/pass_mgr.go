package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type PassManager struct{}

func (p *PassManager) Hash(pass string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (p *PassManager) Check(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(pass), []byte(hash))
}
