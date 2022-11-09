package handlers

import (
	"os"

	"github.com/gorilla/mux"
)

func Handlers() {
	router := mux.NewRouter()

	PORT := os.Getenv("PORT")
}
