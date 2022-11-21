package main

import (
	"db"
	"handlers"
	"log"
)

func main() {
	if !db.DbConn.IsConnection() {
		log.Fatal("No connection to the DB")
	}

	handlers.Handlers()
}
