package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ComparePassword(hash, password string) bool {
	byteHash := []byte(hash)
	bytePwd := []byte(password)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
