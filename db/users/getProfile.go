package users

import (
	"context"
	"db"
	"helpers"
	"log"
	"models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* GetProfile Gets a profile in the DB */
func GetProfile(Id string) (models.User, error) {
	timeout, _ := time.ParseDuration(helpers.GetEnvVariable("DB_TIMEOUT"))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	database := db.MongoConnection.Database("twittor")
	col := database.Collection("users")

	var profile models.User

	objId, _ := primitive.ObjectIDFromHex(Id)
	condition := bson.M{
		"_id": objId,
	}

	err := col.FindOne(ctx, condition).Decode(&profile)

	if err != nil {
		log.Println("Registry not found: " + err.Error())

		return profile, err
	}

	profile.Password = ""

	return profile, nil
}
