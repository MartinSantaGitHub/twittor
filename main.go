package main

import (
	"db"
	"handlers"
	"log"
)

func main() {
	if !db.IsConnection() {
		log.Fatal("No connection to the DB")
	}

	handlers.Handlers()
}
