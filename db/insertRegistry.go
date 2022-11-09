package db

import (
	"context"
	"helpers"
	"models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertRegistry inserts an user into de DB */
func InsertRegistry(u models.User) (string, bool, error) {
	timeout, _ := time.ParseDuration(helpers.GetEnvVariable("DB_TIMEOUT"))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	db := MongoConnection.Database("twittor")
	col := db.Collection("users")

	u.Password, _ = EncryptPassword(u.Password)

	result, err := col.InsertOne(ctx, u)

	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)

	return ObjID.String(), true, nil
}
