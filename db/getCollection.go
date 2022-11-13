package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

/* GetCollection Gets a collection from a mongo DB */
func GetCollection(database string, collection string) *mongo.Collection {
	db := MongoConnection.Database(database)
	col := db.Collection(collection)

	return col
}
