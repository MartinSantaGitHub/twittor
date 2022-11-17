package handlers

import (
	"fmt"
	"helpers"
	"log"
	"net/http"
	"routes/relations"
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

	// Register Home page service
	router.HandleFunc("/", home)

	// Register Users endpoints
	users.Insert(router)
	users.Login(router)
	users.GetProfile(router)
	users.Modify(router)
	users.UploadAvatar(router)
	users.UploadBanner(router)
	users.GetAvatar(router)
	users.GetBanner(router)

	// Register Tweets endpoints
	tweets.Insert(router)
	tweets.GetTweets(router)
	tweets.Delete(router)

	// Register Relations endpoints
	relations.Insert(router)
	relations.Delete(router)
	relations.IsRelation(router)
	relations.GetUsers(router)

	PORT := helpers.GetEnvVariable("PORT")
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler))
}
