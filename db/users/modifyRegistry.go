package users

import (
	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* ModifyRegistry modifies a registry in the DB */
func ModifyRegistry(user models.User, id string) (bool, error) {
	col := db.GetCollection("twittor", "users")
	registry := make(map[string]interface{})

	if len(user.Name) > 0 {
		registry["name"] = user.Name
	}

	if len(user.LastName) > 0 {
		registry["lastName"] = user.LastName
	}

	if len(user.Avatar) > 0 {
		registry["Avatar"] = user.Avatar
	}

	if len(user.Banner) > 0 {
		registry["banner"] = user.Banner
	}

	if len(user.Biography) > 0 {
		registry["biography"] = user.Biography
	}

	if len(user.Location) > 0 {
		registry["location"] = user.Location
	}

	if len(user.WebSite) > 0 {
		registry["webSite"] = user.WebSite
	}

	registry["birthDate"] = user.BirthDate

	updateString := bson.M{
		"$set": registry,
	}

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	//filter := bson.M{"_id": bson.M{"$eq": objId}}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	_, err := col.UpdateOne(ctx, filter, updateString)
	//_, err := col.UpdateByID(ctx, objId, updateString)

	defer cancel()

	if err != nil {
		return false, err
	}

	return true, nil
}
