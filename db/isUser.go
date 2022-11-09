package db

import (
	"context"
	"helpers"
	"models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/* IsUser checks that the user already exists in the DB */
func IsUser(email string) (models.User, bool, string) {
	timeout, _ := time.ParseDuration(helpers.GetEnvVariable("DB_TIMEOUT"))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	db := MongoConnection.Database("twittor")
	col := db.Collection("users")
	condition := bson.M{"email": email}

	var result models.User

	err := col.FindOne(ctx, condition).Decode(&result)
	Id := result.Id.Hex()

	if err != nil {
		return result, false, Id
	}

	return result, true, Id
}
