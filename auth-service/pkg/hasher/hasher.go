package hasher

import (
	"errors"

	"github.com/monstrong/gracker2/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type SimpleHasher struct{
	cost int
}

func NewSimpleHasher(cost int) *SimpleHasher {
	return &SimpleHasher{cost: cost}
}

func (s *SimpleHasher) HashPassword(password string) (string, error) { 
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func (s *SimpleHasher) CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
        return errs.ErrPasswordMismatch
    }
	return err
}

