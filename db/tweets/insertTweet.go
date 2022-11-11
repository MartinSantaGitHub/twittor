package tweets

import (
	"db"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertTweet Insert a tweet in the DB */
func InsertTweet(tweet models.Tweet) (string, bool, error) {
	col, ctx, cancel := db.GetCollection("twittor", "tweet")

	defer cancel()

	registry := bson.M{
		"userId":  tweet.UserId,
		"message": tweet.Message,
		"date":    tweet.Date,
	}

	result, err := col.InsertOne(ctx, registry)

	if err != nil {
		return "", false, err
	}

	objId := result.InsertedID.(primitive.ObjectID)

	// The same goes with objId.Hex()
	return objId.String(), true, nil
}
