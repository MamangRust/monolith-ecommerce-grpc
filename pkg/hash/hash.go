package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

//go:generate mockgen -source=hash.go -destination=mocks/hash.go
type HashPassword interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashPassword string, password string) error
}

type Hashing struct{}

func NewHashingPassword() HashPassword {
	return &Hashing{}
}

func (h Hashing) HashPassword(password string) (string, error) {
	pw := []byte(password)
	hashedPw, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPw), nil
}

func (h Hashing) ComparePassword(hashPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
