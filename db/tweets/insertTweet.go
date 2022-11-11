package tweets

import (
	"context"
	"db"
	"helpers"
	"models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertTweet Insert a tweet in the DB */
func InsertTweet(tweet models.Tweet) (string, bool, error) {

	timeout, _ := time.ParseDuration(helpers.GetEnvVariable("DB_TIMEOUT"))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	database := db.MongoConnection.Database("twittor")
	col := database.Collection("tweet")
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
