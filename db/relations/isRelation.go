package relations

import (
	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/mongo"
)

/* IsRelation verifies if exist a relation in the DB */
func IsRelation(relation models.Relation) (bool, models.Relation, error) {
	var result models.Relation

	col := db.GetCollection("twittor", "relation")
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	err := col.FindOne(ctx, relation).Decode(&result)

	defer cancel()

	if err != nil && err == mongo.ErrNoDocuments {
		return false, result, nil
	}

	return true, result, err
}
