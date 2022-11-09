package routers

import (
	"db"
	"encoding/json"
	"models"
	"net/http"
)

/* Registry permits to create a user in the DB */
func Registry(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "There was an error receiving the data: "+err.Error(), 400)

		return
	}

	if len(user.Email) == 0 {
		http.Error(w, "The email is required", 400)

		return
	}

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
