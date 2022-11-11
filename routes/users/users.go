package users

import (
	"controllers/users"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* Insert Allows to create an user */
func Insert(router *mux.Router) {
	router.HandleFunc("/user", helpers.MultipleMiddleware(users.Insert, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* Login Permits a user to login in the service */
func Login(router *mux.Router) {
	router.HandleFunc("/user/login", helpers.MultipleMiddleware(users.Login, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* GetProfile Gets an user profile */
func GetProfile(router *mux.Router) {
	router.HandleFunc("/user/profile", helpers.MultipleMiddleware(users.GetProfile, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("GET")
}

/* Modify Allows to modify a registry */
func Modify(router *mux.Router) {
	router.HandleFunc("/user", helpers.MultipleMiddleware(users.Modify, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("PUT")
}
