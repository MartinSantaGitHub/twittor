package users

import (
	"db"
	"log"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* GetProfile Gets a profile in the DB */
func GetProfile(Id string) (models.User, error) {
	col, ctx, cancel := db.GetCollection("twittor", "users")

	defer cancel()

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
