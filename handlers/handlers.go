package handlers

import (
	"fmt"
	"helpers"
	"log"
	"net/http"
	"routes/tweets"
	"routes/users"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

/* Handler that set the PORT and run the service */
func Handlers() {
	router := mux.NewRouter()

	router.HandleFunc("/", home)

	// Register Users endpoints
	users.Insert(router)
	users.Login(router)
	users.GetProfile(router)
	users.Modify(router)

	// Register Tweets endpoints
	tweets.Insert(router)
	tweets.Get(router)

	PORT := helpers.GetEnvVariable("PORT")
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler))
}
