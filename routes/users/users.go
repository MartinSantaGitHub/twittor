package routes

import (
	"controllers/users"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

func Registry(router *mux.Router) {
	router.HandleFunc("/registry", helpers.MultipleMiddleware(users.Registry, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}
