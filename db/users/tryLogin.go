package users

import (
	"models"

	"golang.org/x/crypto/bcrypt"
)

/* TryLogin makes the login to the DB */
func TryLogin(email string, password string) (models.User, bool) {
	user, isFound, _ := IsUser(email)

	if !isFound {
		return user, false
	}

	passwordBytes := []byte(password)
	passwordDB := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)

	if err != nil {
		return user, false
	}

	return user, true
}
