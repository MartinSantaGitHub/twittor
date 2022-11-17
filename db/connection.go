package db

import (
	"context"
	"fmt"
	"helpers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(helpers.GetEnvVariable("MONGO_CONN"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println(err.Error())

		return client
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err.Error())

		return client
	}

	fmt.Println("Connection successful to the DB")

	return client
}

/* MongoConnection is the connection object to the Database */
var MongoConnection = connectDB()

/* CheckConnection makes a ping to the Database */
func IsConnection() bool {
	err := MongoConnection.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err.Error())

		return false
	}

	return true
}
