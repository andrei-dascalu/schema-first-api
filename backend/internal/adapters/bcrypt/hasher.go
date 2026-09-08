// Package bcrypt implements app.PasswordHasher with bcrypt.
package bcrypt

import (
	"golang.org/x/crypto/bcrypt"

	"spec-first-backend/internal/app"
)

var _ app.PasswordHasher = (*Hasher)(nil)

type Hasher struct{}

func New() *Hasher {
	return &Hasher{}
}

func (*Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (*Hasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
