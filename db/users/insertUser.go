package users

import (
	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertUser inserts an user into de DB */
func InsertUser(u models.User) (string, bool, error) {
	col := db.GetCollection("twittor", "users")

	u.Password, _ = EncryptPassword(u.Password)

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	result, err := col.InsertOne(ctx, u)

	defer cancel()

	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)

	return ObjID.String(), true, nil
}
