package users

import "golang.org/x/crypto/bcrypt"

/* EncryptPassword encrypts the user's password */
func EncryptPassword(password string) (string, error) {
	// Minimum - cost: 6
	// Common user - cost: 6
	// Admin user - cost: 8

	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}
