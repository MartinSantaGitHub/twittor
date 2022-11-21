package middlewares

import (
	"db"
	"net/http"
)

/* CheckDB is the middleware that allows to know the status of the DB */
func CheckDB(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !db.DbConn.IsConnection() {
			http.Error(w, "Connection with the DB lost", 500)

			return
		}

		next.ServeHTTP(w, r)
	}
}
