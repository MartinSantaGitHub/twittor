package users

import (
	"db"
	"helpers"
	"models"
	"net/http"
	"net/mail"
)

/* Registry permits to create a user in the DB */
func Registry(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(helpers.RequestUserKey{}).(models.User)

	if len(user.Password) < 6 {
		http.Error(w, "The password must have at least 6 characters", 400)

		return
	}

	_, isFound, _ := db.IsUser(user.Email)

	if isFound {
		http.Error(w, "The user already exists", 400)

		return
	}

	_, status, err := db.InsertRegistry(user)

	if err != nil {
		http.Error(w, "There was an error trying to regist the user"+err.Error(), 400)

		return
	}

	if !status {
		http.Error(w, "The user registry could not be inserted into the DB", 400)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* Login does the login */
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	user := r.Context().Value(helpers.RequestUserKey{}).(models.User)

	if len(user.Email) == 0 {
		http.Error(w, "The email is required", 400)

		return
	}

	_, err := mail.ParseAddress(user.Email)

	if err != nil {
		http.Error(w, "Invalid email format", 400)

		return
	}
}
