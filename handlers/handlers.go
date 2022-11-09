package handlers

import (
	"fmt"
	"helpers"
	"log"
	"middlewares"
	"net/http"
	"routers"

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
	router.HandleFunc("/registry", middlewares.CheckDB(routers.Registry)).Methods("POST")

	PORT := helpers.GetEnvVariable("PORT")
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler))
}
