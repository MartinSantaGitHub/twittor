package users

import (
	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
)

/* IsUser checks that the user already exists in the DB */
func IsUser(email string) (models.User, bool, string) {
	var result models.User

	col := db.GetCollection("twittor", "users")
	condition := bson.M{"email": email}
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	err := col.FindOne(ctx, condition).Decode(&result)
	id := result.Id.Hex()

	defer cancel()

	if err != nil {
		return result, false, id
	}

	return result, true, id
}
