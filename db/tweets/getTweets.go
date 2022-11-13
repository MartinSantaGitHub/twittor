package tweets

import (
	"context"
	"db"
	"log"

	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Get Get a user's tweets from the DB */
func GetTweets(id string, page int64, limit int64) ([]*models.Tweet, int64, bool) {
	ctxCount, cancelCount := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	col, ctx, cancel := db.GetCollection("twittor", "tweet")

	defer func() {
		cancelCount()
		cancel()
	}()

	var results []*models.Tweet

	condition := bson.M{
		"userId": id,
		"active": true,
	}

	total, err := col.CountDocuments(ctxCount, condition)

	if err != nil {
		log.Fatal(err.Error())

		return results, total, false
	}

	opts := options.Find()

	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * limit)

	cursor, err := col.Find(ctx, condition, opts)

	if err != nil {
		log.Fatal(err.Error())

		return results, total, false
	}

	for cursor.Next(context.TODO()) {
		var registry models.Tweet

		err := cursor.Decode(&registry)

		if err != nil {
			return results, total, false
		}

		results = append(results, &registry)
	}

	return results, total, true
}
