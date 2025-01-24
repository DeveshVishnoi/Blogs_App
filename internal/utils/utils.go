package utils

import (
	"social/internal/store"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(text string) (store.Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return store.Password{}, err
	}

	password := store.Password{
		Text: text,
		Hash: hash,
	}

	return password, nil

}

func Compare(text string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(text))
}
