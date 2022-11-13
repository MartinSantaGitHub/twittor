package relations

import (
	"db"
	"helpers"
	"models"
)

/* Delete Deletes a relation in the DB */
func DeleteFisical(relation models.Relation) error {
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	col := db.GetCollection("twittor", "relation")
	_, err := col.DeleteOne(ctx, relation)

	defer cancel()

	return err
}

/* DeleteLogical Inactivates a relation in the DB */
func DeleteLogical(relation models.Relation) error {
	col := db.GetCollection("twittor", "relation")
	updateString := map[string]map[string]bool{"$set": {"active": false}}

	// Also bson.M{"$set": bson.M{"active": false},} in the updateString

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	_, err := col.UpdateOne(ctx, relation, updateString)

	defer cancel()

	return err
}
