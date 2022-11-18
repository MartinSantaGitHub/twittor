package tweets

import (
	"context"
	"db"

	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Get Get a user's tweets from the DB */
func GetTweets(id primitive.ObjectID, page int64, limit int64) ([]*models.Tweet, int64, error) {
	var results []*models.Tweet

	col := db.GetCollection("twittor", "tweet")
	condition := bson.M{
		"userId": id,
		"active": true,
	}

	ctxCount, cancelCount := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	total, err := col.CountDocuments(ctxCount, condition)

	defer cancelCount()

	if err != nil {
		return results, total, err
	}

	opts := options.Find()

	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	cursor, err := col.Find(ctxFind, condition, opts)

	defer cancelFind()

	if err != nil {
		return results, total, err
	}

	ctxCursor := context.TODO()

	defer cursor.Close(ctxCursor)

	err = cursor.All(ctxCursor, &results)

	if err != nil {
		return results, total, err
	}

	return results, total, nil
}
