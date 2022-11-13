package db

import (
	"context"

	"helpers"

	"go.mongodb.org/mongo-driver/mongo"
)

/* GetCollection Gets a collection from a mongo DB */
func GetCollection(database string, collection string) (*mongo.Collection, context.Context, context.CancelFunc) {
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	db := MongoConnection.Database(database)
	col := db.Collection(collection)

	return col, ctx, cancel
}
