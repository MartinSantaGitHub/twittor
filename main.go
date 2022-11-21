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

	log.Println("Connection successful to the DB")

	handlers.Handlers()
}
