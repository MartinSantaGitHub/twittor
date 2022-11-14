package relations

import (
	"fmt"

	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
)

/* InsertRelation Creates a relation into the DB */
func InsertRelation(relation models.Relation) error {
	var relationDb models.Relation
	var err error

	col := db.GetCollection("twittor", "relation")
	isRelation, relationDb, err := IsRelation(relation)

	if err != nil {
		return err
	}

	if !isRelation {
		ctxInsert, cancelInsert := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
		_, err = col.InsertOne(ctxInsert, relation)

		defer cancelInsert()

		return err
	}

	if relationDb.Active {
		return fmt.Errorf("the relation with the user id: %s already exists", relation.UserRelationId)
	}

	updateString := bson.M{
		"$set": bson.M{"active": true},
	}

	ctxUpdate, cancelUpdate := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	_, err = col.UpdateByID(ctxUpdate, relationDb.Id, updateString)

	defer cancelUpdate()

	return err
}
