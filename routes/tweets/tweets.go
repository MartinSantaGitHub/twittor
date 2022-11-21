package tweets

import (
	"controllers/tweets"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* Insert allows to create a new tweet */
func Insert(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.Insert,
		middlewares.CheckDB,
		middlewares.ValidateJWT)).Methods("POST")
}

/* GetTweets gets an user's tweets */
func GetTweets(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.GetTweets,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId,
		middlewares.ValidatePageLimit)).Methods("GET")
}

/* Delete deletes an user's tweet */
func Delete(router *mux.Router) {
	router.HandleFunc("/tweet", helpers.MultipleMiddleware(tweets.Delete,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("DELETE")
}
