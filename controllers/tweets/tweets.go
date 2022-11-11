package tweets

import (
	"encoding/json"
	"helpers"
	"jwt"
	"net/http"
	"strconv"
	"time"

	db "db/tweets"
	"models"
	mr "models/response"
)

/* InsertTweet permits to insert a tweet in the DB */
func Insert(w http.ResponseWriter, r *http.Request) {
	var tweet models.Tweet

	err := json.NewDecoder(r.Body).Decode(&tweet)

	if err != nil {
		http.Error(w, "Invalid data"+err.Error(), http.StatusBadRequest)

		return
	}

	if len(tweet.Message) < 1 {
		http.Error(w, "The message cannot be empty", http.StatusBadRequest)

		return
	}

	registry := models.Tweet{
		UserId:  jwt.UserId,
		Message: tweet.Message,
		Date:    time.Now(),
	}

	_, status, err := db.InsertTweet(registry)

	if err != nil {
		http.Error(w, "An error occurred trying to insert a new registry into the DB: "+err.Error(), http.StatusInternalServerError)

		return
	}

	if !status {
		http.Error(w, "Could not be inserted a registry into the DB", http.StatusNotModified)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* GetTweets GetTweets a user's tweets */
func GetTweets(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if len(id) < 1 {
		http.Error(w, "The id param is required", http.StatusBadRequest)

		return
	}

	pageQuery := r.URL.Query().Get("page")

	if len(pageQuery) < 1 {
		pageQuery = "1"
	}

	limitQuery := r.URL.Query().Get("limit")

	if len(limitQuery) < 1 {
		limitQuery = helpers.GetEnvVariable("TWEETS_GET_LIMIT")
	}

	page, err := strconv.ParseInt(pageQuery, 10, 64)

	if err != nil {
		http.Error(w, "The page param is invalid", http.StatusBadRequest)

		return
	}

	limit, err := strconv.ParseInt(limitQuery, 10, 64)

	if err != nil {
		http.Error(w, "The limit param is invalid", http.StatusBadRequest)

		return
	}

	results, total, success := db.GetTweets(id, page, limit)

	if !success {
		http.Error(w, "An error has happened trying to get the tweets from the DB "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := mr.TweetsResponse{
		Tweets: results,
		Total:  total,
	}

	json.NewEncoder(w).Encode(response)
}

/* Delete Deletes a tweet that belongs to an user */
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if len(id) < 1 {
		http.Error(w, "The id param is required", http.StatusBadRequest)

		return
	}

	err := db.DeleteLogic(id, jwt.UserId)

	if err != nil {
		http.Error(w, "An error occurred trying to delete the tweet "+err.Error(), http.StatusInternalServerError)

		return
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
