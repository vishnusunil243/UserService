package helper

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("error hashing password")
		return ""
	}
	return string(hash)
}
