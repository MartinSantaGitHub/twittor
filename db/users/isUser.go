package users

import (
	"db"
	"models"

	"go.mongodb.org/mongo-driver/bson"
)

/* IsUser checks that the user already exists in the DB */
func IsUser(email string) (models.User, bool, string) {
	col, ctx, cancel := db.GetCollection("twittor", "users")

	defer cancel()

	condition := bson.M{"email": email}

	var result models.User

	err := col.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()

	if err != nil {
		return result, false, id
	}

	return result, true, id
}
