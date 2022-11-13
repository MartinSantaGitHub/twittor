package tweets

import (
	"db"
	"jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Delete Deletes a tweet in the DB */
func DeleteFisical(id string, userId string) error {
	col, ctx, cancel := db.GetCollection("twittor", "tweet")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{
		"_id":    objId,
		"userId": jwt.UserId,
	}

	_, err := col.DeleteOne(ctx, condition)

	return err
}

/* DeleteLogic Inactivates a tweet in the DB */
func DeleteLogic(id string, userId string) error {
	col, ctx, cancel := db.GetCollection("twittor", "tweet")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{
		"_id":    objId,
		"userId": jwt.UserId,
	}
	updateString := bson.M{
		"$set": bson.M{"active": false},
	}

	// Also map[string]bool{"active": false} in the updateString

	_, err := col.UpdateOne(ctx, condition, updateString)

	return err
}
