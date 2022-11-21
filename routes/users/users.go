package users

import (
	"controllers/users"
	"helpers"
	"middlewares"

	"github.com/gorilla/mux"
)

/* Insert allows to create an user */
func Insert(router *mux.Router) {
	router.HandleFunc("/user", helpers.MultipleMiddleware(users.Insert, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* Login permits an user to login in the service */
func Login(router *mux.Router) {
	router.HandleFunc("/user/login", helpers.MultipleMiddleware(users.Login, middlewares.CheckDB, middlewares.ValidateEmail)).Methods("POST")
}

/* GetProfile gets an user profile */
func GetProfile(router *mux.Router) {
	router.HandleFunc("/user/profile", helpers.MultipleMiddleware(users.GetProfile,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("GET")
}

/* Modify allows to modify a registry */
func Modify(router *mux.Router) {
	router.HandleFunc("/user", helpers.MultipleMiddleware(users.Modify, middlewares.CheckDB, middlewares.ValidateJWT)).Methods("PUT")
}

/* Upload uploads an user's avatar */
func UploadAvatar(router *mux.Router) {
	router.HandleFunc("/user/avatar", helpers.MultipleMiddleware(users.UploadAvatar,
		middlewares.CheckDB,
		middlewares.ValidateJWT)).Methods("POST")
}

/* Upload uploads an user's avatar */
func UploadBanner(router *mux.Router) {
	router.HandleFunc("/user/banner", helpers.MultipleMiddleware(users.UploadBanner,
		middlewares.CheckDB,
		middlewares.ValidateJWT)).Methods("POST")
}

/* GetAvatar gets the user's avatar */
func GetAvatar(router *mux.Router) {
	router.HandleFunc("/user/avatar", helpers.MultipleMiddleware(users.GetAvatar,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("GET")
}

/* GetBanner gets the user's banner */
func GetBanner(router *mux.Router) {
	router.HandleFunc("/user/banner", helpers.MultipleMiddleware(users.GetBanner,
		middlewares.CheckDB,
		middlewares.ValidateJWT,
		middlewares.ValidateQueryId)).Methods("GET")
}
