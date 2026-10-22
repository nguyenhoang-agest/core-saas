package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12

// HashPassword hashes a plain-text password using bcrypt with DefaultCost.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns true if the password matches the hash.
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
