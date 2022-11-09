package handlers

import (
	"fmt"
	"helpers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

/* Handler that set the PORT and run the service */
func Handlers() {
	PORT := helpers.GetEnvVariable("PORT")
	router := mux.NewRouter()
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler))
}
