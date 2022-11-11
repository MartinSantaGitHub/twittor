package users

import (
	"db"
	"models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertUser inserts an user into de DB */
func InsertUser(u models.User) (string, bool, error) {
	col, ctx, cancel := db.GetCollection("twittor", "users")

	defer cancel()

	u.Password, _ = EncryptPassword(u.Password)

	result, err := col.InsertOne(ctx, u)

	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)

	return ObjID.String(), true, nil
}
