package tweets

import (
	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertTweet Insert a tweet in the DB */
func InsertTweet(tweet models.Tweet) (string, bool, error) {
	col := db.GetCollection("twittor", "tweet")
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	result, err := col.InsertOne(ctx, tweet)

	defer cancel()

	if err != nil {
		return "", false, err
	}

	objId := result.InsertedID.(primitive.ObjectID)

	// The same goes with objId.Hex()
	return objId.String(), true, nil
}
