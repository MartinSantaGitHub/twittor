package middlewares

import (
	"context"
	"encoding/json"
	"helpers"
	"models"
	"net/http"
	"net/mail"
)

func ValidateEmail(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err = mail.ParseAddress(user.Email)

		if err != nil {
			http.Error(w, "Invalid email format", 400)

			return
		}

		ctx := context.WithValue(r.Context(), helpers.RequestUserKey{}, user)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}
