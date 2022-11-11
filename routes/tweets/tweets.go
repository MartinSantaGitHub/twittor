package tweets

import (
	"controllers/tweets"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* Insert Allows to create a new tweet */
func Insert(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.Insert, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("POST")
}

/* Get Gets an user's tweets */
func Get(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.Get, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("GET")
}
