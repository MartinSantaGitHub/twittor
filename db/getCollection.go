package db

import (
	"context"
	"time"

	"helpers"

	"go.mongodb.org/mongo-driver/mongo"
)

/* GetCollection Gets a collection from a mongo DB */
func GetCollection(database string, collection string) (*mongo.Collection, context.Context, context.CancelFunc) {
	timeout, _ := time.ParseDuration(helpers.GetEnvVariable("DB_TIMEOUT"))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	db := MongoConnection.Database(database)
	col := db.Collection(collection)

	return col, ctx, cancel
}
