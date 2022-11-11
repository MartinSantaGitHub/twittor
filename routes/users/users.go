package routes

import (
	"controllers/users"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* Registry Allows to create an user */
func Registry(router *mux.Router) {
	router.HandleFunc("/registry", helpers.MultipleMiddleware(users.Registry, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* Login Permits a user to login in the service */
func Login(router *mux.Router) {
	router.HandleFunc("/login", helpers.MultipleMiddleware(users.Login, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* GetProfile Gets an user profile */
func GetProfile(router *mux.Router) {
	router.HandleFunc("/profile", helpers.MultipleMiddleware(users.GetProfile, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("GET")
}

/* ModifyRegistry Allows to modify a registry */
func ModifyRegistry(router *mux.Router) {
	router.HandleFunc("/registry", helpers.MultipleMiddleware(users.ModifyRegistry, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("PUT")
}
