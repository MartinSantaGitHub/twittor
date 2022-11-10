package users

import (
	"db"
	"encoding/json"
	"helpers"
	"jwt"
	"models"
	mr "models/response"
	"net/http"
	"time"
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
	user := r.Context().Value(helpers.RequestUserKey{}).(models.User)
	userDb, isUser := db.TryLogin(user.Email, user.Password)

	if !isUser {
		http.Error(w, "User and/or password invalid", 400)

		return
	}

	jwtKey, err := jwt.GenerateJWT(userDb)

	if err != nil {
		http.Error(w, "Something went wrong"+err.Error(), 500)

		return
	}

	response := mr.LoginResponse{
		Token: jwtKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

	expirationTime := time.Now().Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: expirationTime,
	})
}
