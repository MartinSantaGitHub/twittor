package bd

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions = options.Client().ApplyURI("mongodb+srv://admin_clus_cafe:Capgemini006@cluscafe.taxbz4x.mongodb.net/?retryWrites=true&w=majority")

/* MongoConnection is the connection object to the Database */
var MongoConnection = ConnectDB()

/* ConnectDB allows to connect to the Database */
func ConnectDB() *mongo.Client {
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

/* CheckConnection makes a ping to the Database */
func CheckConnection() int {
	err := MongoConnection.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err.Error())

		return 0
	}

	return 1
}
