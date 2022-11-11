package tweets

import (
	"encoding/json"
	"jwt"
	"net/http"
	"time"

	db "db/tweets"
	"models"
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
