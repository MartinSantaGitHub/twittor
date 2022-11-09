package db

import (
	"context"
	"helpers"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(helpers.GetEnvVariable("MONGO_CONN"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err.Error())

		return client
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err.Error())

		return client
	}

	log.Println("Connection successful to the DB")

	return client
}

/* MongoConnection is the connection object to the Database */
var MongoConnection = connectDB()

/* CheckConnection makes a ping to the Database */
func IsConnection() bool {
	err := MongoConnection.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err.Error())

		return false
	}

	return true
}
