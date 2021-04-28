package hash

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) (*SHA1Hasher, error) {
	if salt == "" {
		return nil, errors.New("empty hasher salt")
	}
	return &SHA1Hasher{salt: salt}, nil
}

func (h *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
