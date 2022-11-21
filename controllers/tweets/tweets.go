package tweets

import (
	"encoding/json"
	"net/http"
	"time"

	"db"
	"helpers"
	"jwt"
	req "models/request"
	res "models/response"
)

/* InsertTweet creates a tweet in the DB */
func Insert(w http.ResponseWriter, r *http.Request) {
	var tweet req.Tweet

	err := json.NewDecoder(r.Body).Decode(&tweet)

	if err != nil {
		http.Error(w, "Invalid data: "+err.Error(), http.StatusBadRequest)

		return
	}

	if len(tweet.Message) < 1 {
		http.Error(w, "The message cannot be empty", http.StatusBadRequest)

		return
	}

	registry := req.Tweet{
		UserId:  jwt.UserId,
		Message: tweet.Message,
		Date:    time.Now(),
		Active:  true,
	}

	_, err = db.DbConn.InsertTweet(registry)

	if err != nil {
		http.Error(w, "An error occurred trying to insert a new registry into the DB: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* GetTweets gets an user's tweets */
func GetTweets(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	page := r.Context().Value(helpers.RequestPageKey{}).(int64)
	limit := r.Context().Value(helpers.RequestLimitKey{}).(int64)

	results, total, err := db.DbConn.GetTweets(id, page, limit)

	if err != nil {
		http.Error(w, "An error has happened trying to get the tweets from the DB "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := res.TweetsResponse{
		Tweets: results,
		Total:  total,
	}

	json.NewEncoder(w).Encode(response)
}

/* Delete deletes a tweet that belongs to an user */
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	err := db.DbConn.DeleteTweetLogical(id)

	if err != nil {
		http.Error(w, "An error occurred trying to delete the tweet: "+err.Error(), http.StatusInternalServerError)

		return
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
