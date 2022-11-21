package middlewares

import (
	"jwt"
	"net/http"
)

/* ValidateJWT allows to validate the JWT from the request */
func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.ProcessJWT(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Error on the Token: "+err.Error(), http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	}
}
