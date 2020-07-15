package bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PassManager struct{}

func (p *PassManager) Hash(pass string) (string, error) {
	fmt.Println("TRYING TO GENERATE HASH")
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	fmt.Println("NEW HASH", string(h))
	return string(h), nil
}

func (p *PassManager) Check(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(pass), []byte(hash))
}
