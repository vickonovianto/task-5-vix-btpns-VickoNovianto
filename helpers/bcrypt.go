package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainTextPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func IsPasswordMatch(hashedPassword string, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
}
