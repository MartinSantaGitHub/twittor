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

/* GetTweets Gets an user's tweets */
func GetTweets(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.GetTweets, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("GET")
}

/* Delete Deletes an user's tweet */
func Delete(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.Delete, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("DELETE")
}
