package tweets

import (
	"db"
	"helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Delete Deletes a tweet in the DB */
func DeleteFisical(id string) error {
	col := db.GetCollection("twittor", "tweet")
	objId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{
		"_id": objId,
	}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	_, err := col.DeleteOne(ctx, condition)

	defer cancel()

	return err
}

/* DeleteLogical Inactivates a tweet in the DB */
func DeleteLogical(id string) error {
	col := db.GetCollection("twittor", "tweet")
	objId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{
		"_id": objId,
	}
	updateString := bson.M{
		"$set": bson.M{"active": false},
	}

	// Also map[string]map[string]bool{"$set": {"active": false}} in the updateString

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	_, err := col.UpdateOne(ctx, condition, updateString)

	defer cancel()

	return err
}
