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
		log.Fatal(err.Error())

		return results, total, false
	}

	opts := options.Find()

	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	cursor, err := col.Find(ctxFind, condition, opts)

	defer cancelFind()

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
