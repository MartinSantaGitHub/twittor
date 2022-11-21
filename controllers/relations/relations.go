package relations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"db"
	"helpers"
	"jwt"
	req "models/request"
	res "models/response"
)

// region "Actions"

/* Creates creates a new relation between two users */
func Create(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)

	if jwt.UserId == id {
		http.Error(w, fmt.Sprintf("Error: %s", "relation with oneself not allowed (userId and userRelationId are the same)"), http.StatusBadRequest)

		return
	}

	relation, err := getRelationRequestModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = db.DbConn.InsertRelation(relation)

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

/* Delete deletes a relation */
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation, err := getRelationRequestModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = db.DbConn.DeleteRelationLogical(relation)

	if err != nil {
		http.Error(w, fmt.Sprintf("An error has occurred trying to delete a relation: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* IsRelation checks if exist a relation */
func IsRelation(w http.ResponseWriter, r *http.Request) {
	var isRelation bool

	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	relation, err := getRelationRequestModel(id)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)

		return
	}

	found, relationDb, err := db.DbConn.GetRelation(relation)

	if err != nil {
		http.Error(w, fmt.Sprintf("An error has occurred trying to obtain a relation: %s", err.Error()), http.StatusInternalServerError)

		return
	}

	if found {
		isRelation = relationDb.Active
	}

	response := res.IsRelationResponse{
		Status: isRelation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/* GetUsers gets a list of users */
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var results []*req.User
	var total int64
	var err error

	page := r.Context().Value(helpers.RequestPageKey{}).(int64)
	limit := r.Context().Value(helpers.RequestLimitKey{}).(int64)
	search := r.URL.Query().Get("search")
	searchType := r.URL.Query().Get("type")

	switch searchType {

	case "new":
		results, total, err = db.DbConn.GetNotFollowers(jwt.UserId, page, limit, search)
	case "follow":
		results, total, err = db.DbConn.GetFollowers(jwt.UserId, page, limit, search)
	default:
		http.Error(w, "Invalid type param value. It has to be \"follow\" or \"new\"", http.StatusBadRequest)

		return
	}

	//results, total, err = db.DbConn.GetUsers(jwt.UserId, page, limit, search, searchType)

	if err != nil {
		http.Error(w, "Error getting the users: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := res.UsersResponse{
		Users: results,
		Total: total,
	}

	json.NewEncoder(w).Encode(response)
}

/* GetUsersTweets returns the followers' tweets */
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

	results, total, err := db.DbConn.GetUsersTweets(jwt.UserId, page, limit, isOnlyTweets)

	if err != nil {
		http.Error(w, "Error getting the tweets: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	/* This is optional, the default status value is 200 OK */
	//w.WriteHeader(http.StatusOK)

	if isOnlyTweets {
		response = res.TweetsResponse{
			Tweets: results.([]*req.Tweet),
			Total:  total,
		}
	} else {
		response = res.UserTweetsResponse{
			Tweets: results.([]*req.UserTweet),
			Total:  total,
		}
	}

	json.NewEncoder(w).Encode(response)
}

// endregion

// region "Helpers"

func getRelationRequestModel(userRelationId string) (req.Relation, error) {
	var relation req.Relation

	relation.UserId = jwt.UserId
	relation.UserRelationId = userRelationId
	relation.Active = true

	return relation, nil
}

// endregion
