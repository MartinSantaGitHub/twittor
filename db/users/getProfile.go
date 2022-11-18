package users

import (
	"db"
	"helpers"
	"log"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* GetProfile Gets a profile in the DB */
func GetProfile(id primitive.ObjectID) (models.User, error) {
	var profile models.User

	col := db.GetCollection("twittor", "users")

	condition := bson.M{
		"_id": id,
	}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	err := col.FindOne(ctx, condition).Decode(&profile)

	defer cancel()

	if err != nil {
		log.Println("Registry not found: " + err.Error())

		return profile, err
	}

	profile.Password = ""

	return profile, nil
}
