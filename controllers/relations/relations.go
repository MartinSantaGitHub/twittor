package relations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "db/relations"
	"helpers"
	"jwt"
	"models"
	mr "models/response"
	mres "models/result"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Creates a new relation between two users */
func Create(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(primitive.ObjectID)
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
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(primitive.ObjectID)
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
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(primitive.ObjectID)
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

	id, _ := primitive.ObjectIDFromHex(jwt.UserId)

	switch searchType {

	case "new":
		results, total, err = db.GetNotFollowers(id, page, limit, search)
	case "follow":
		results, total, err = db.GetFollowers(id, page, limit, search)
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

	response := mr.UsersResponse{
		Users: results,
		Total: total,
	}

	json.NewEncoder(w).Encode(response)
}

/* GetUsersTweets Returns the followers' tweets */
func GetUsersTweets(w http.ResponseWriter, r *http.Request) {
	var response interface{}

	page := r.Context().Value(helpers.RequestPageKey{}).(int64)
	limit := r.Context().Value(helpers.RequestLimitKey{}).(int64)
	onlyTweets := r.URL.Query().Get("onlytweets")

	if len(onlyTweets) < 1 {
		onlyTweets = "false"
	}

	isOnlyTweets, err := strconv.ParseBool(onlyTweets)

	if err != nil {
		http.Error(w, "Invalid onlytweets param value. It has to be a boolean value", http.StatusBadRequest)

		return
	}

	id, _ := primitive.ObjectIDFromHex(jwt.UserId)

	results, total, err := db.GetUsersTweets(id, page, limit, isOnlyTweets)

	if err != nil {
		http.Error(w, "Error getting the tweets: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	/* This is optional, the default status value is 200 OK */
	//w.WriteHeader(http.StatusOK)

	if isOnlyTweets {
		response = mr.TweetsResponse{
			Tweets: results.([]*models.Tweet),
			Total:  total,
		}
	} else {
		response = mr.UserTweetsResponse{
			Tweets: results.([]*mres.UserTweet),
			Total:  total,
		}
	}

	json.NewEncoder(w).Encode(response)
}

func getRelationModel(userRelationId primitive.ObjectID) (models.Relation, error) {
	var relation models.Relation

	userId, _ := primitive.ObjectIDFromHex(jwt.UserId)

	if userId == userRelationId {
		return relation, fmt.Errorf("relation with oneself not allowed (userId and userRelationId are the same)")
	}

	relation.UserId = userId
	relation.UserRelationId = userRelationId
	relation.Active = true

	return relation, nil
}
