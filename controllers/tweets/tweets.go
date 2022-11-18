package tweets

import (
	"encoding/json"
	"helpers"
	"jwt"
	"net/http"
	"time"

	db "db/tweets"
	"models"
	mr "models/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	objId, _ := primitive.ObjectIDFromHex(jwt.UserId)

	registry := models.Tweet{
		UserId:  objId,
		Message: tweet.Message,
		Date:    time.Now(),
		Active:  true,
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

/* GetTweets Gets a user's tweets */
func GetTweets(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(primitive.ObjectID)
	page := r.Context().Value(helpers.RequestPageKey{}).(int64)
	limit := r.Context().Value(helpers.RequestLimitKey{}).(int64)

	results, total, err := db.GetTweets(id, page, limit)

	if err != nil {
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
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(primitive.ObjectID)
	err := db.DeleteLogical(id)

	if err != nil {
		http.Error(w, "An error occurred trying to delete the tweet "+err.Error(), http.StatusInternalServerError)

		return
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
