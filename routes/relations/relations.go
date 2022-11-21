package relations

import (
	"controllers/relations"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* InsertRelation inserts a new relation between two users */
func Insert(router *mux.Router) {
	router.HandleFunc("/relation", helpers.MultipleMiddleware(relations.Create,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("POST")
}

/* Delete deletes a relation */
func Delete(router *mux.Router) {
	router.HandleFunc("/relation", helpers.MultipleMiddleware(relations.Delete,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("DELETE")
}

/* IsRelation checks if exist a relation */
func IsRelation(router *mux.Router) {
	router.HandleFunc("/relation", helpers.MultipleMiddleware(relations.IsRelation,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("GET")
}

/* GetUsers gets a list of users */
func GetUsers(router *mux.Router) {
	router.HandleFunc("/relation/users", helpers.MultipleMiddleware(relations.GetUsers,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidatePageLimit)).Methods("GET")
}

/* GetUsersTweets returns the followers' tweets */
func GetUsersTweets(router *mux.Router) {
	router.HandleFunc("/relation/tweets", helpers.MultipleMiddleware(relations.GetUsersTweets,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidatePageLimit)).Methods("GET")
}
