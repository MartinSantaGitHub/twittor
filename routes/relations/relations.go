package relations

import (
	"controllers/relations"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* InsertRelation Inserts a new relation between two users */
func Insert(router *mux.Router) {
	router.HandleFunc("/relation", helpers.MultipleMiddleware(relations.Create,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("POST")
}

func Delete(router *mux.Router) {
	router.HandleFunc("/relation", helpers.MultipleMiddleware(relations.Delete,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("DELETE")
}
