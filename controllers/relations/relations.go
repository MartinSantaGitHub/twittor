package relations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "db/relations"
	"helpers"
	"jwt"
	"models"
	mr "models/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Creates a new relation between two users */
func Create(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation, err := getRelationModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = db.InsertRelation(relation)

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
	relation, err := getRelationModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = db.DeleteLogical(relation)

	if err != nil {
		http.Error(w, fmt.Sprintf("An error has occurred trying to delete a relation: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* IsRelation check if exist a relation */
func IsRelation(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation, err := getRelationModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	isRelation, _, err := db.IsRelation(relation)

	if err != nil {
		http.Error(w, fmt.Sprintf("An error has occurred trying to obtain a relation: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	response := mr.IsRelationResponse{
		Status: isRelation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/* GetUsers Gets a list of users */
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var results []*models.User
	var total int64
	var err error

	page := r.Context().Value(helpers.RequestPageKey{}).(int64)
	limit := r.Context().Value(helpers.RequestLimitKey{}).(int64)
	search := r.URL.Query().Get("search")
	searchType := r.URL.Query().Get("type")

	switch searchType {

	case "new":
		results, total, err = db.GetNotFollowers(jwt.UserId, page, limit, search)
	case "follow":
		results, total, err = db.GetFollowers(jwt.UserId, page, limit, search)
	default:
		http.Error(w, "Invalid type param value. It has to be \"follow\" or \"new\"", http.StatusBadRequest)

		return
	}

	//results, total, err = db.GetUsers(jwt.UserId, page, limit, search, searchType)

	if err != nil {
		http.Error(w, "Error getting the users: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := mr.RelationUsersResponse{
		Users: results,
		Total: total,
	}

	json.NewEncoder(w).Encode(response)
}

func getRelationModel(userRelationId string) (models.Relation, error) {
	var relation models.Relation

	if jwt.UserId == userRelationId {
		return relation, fmt.Errorf("relation with oneself not allowed (userId and userRelationId are the same)")
	}

	userId, _ := primitive.ObjectIDFromHex(jwt.UserId)
	uRelationId, err := primitive.ObjectIDFromHex(userRelationId)

	if err != nil {
		return relation, fmt.Errorf("userRelationId: not a valid mongo id format")
	}

	relation.UserId = userId
	relation.UserRelationId = uRelationId
	relation.Active = true

	return relation, nil
}
