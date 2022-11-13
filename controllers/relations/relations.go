package relations

import (
	"fmt"
	"net/http"
	"strings"

	"helpers"
	"jwt"
	"models"

	db "db/relations"
)

/* Creates a new relation between two users */
func Create(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation := getRelationModel(id)
	err := db.InsertRelation(relation)

	if err != nil {
		statusCode := http.StatusInternalServerError
		errStr := err.Error()

		if strings.Contains(errStr, "already exists") {
			statusCode = http.StatusBadRequest
		}

		http.Error(w, fmt.Sprintf("An error has occurred trying to insert a new relation: %s", errStr), statusCode)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* Delete Deletes a relation */
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation := getRelationModel(id)
	err := db.DeleteLogical(relation)

	if err != nil {
		http.Error(w, fmt.Sprintf("An error has occurred trying to delete a relation: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getRelationModel(userRelationId string) models.Relation {
	var relation models.Relation

	relation.UserId = jwt.UserId
	relation.UserRelationId = userRelationId
	relation.Active = true

	return relation
}
