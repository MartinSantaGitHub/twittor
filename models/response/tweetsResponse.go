package response

import "models"

/* TweetsResponse Response model for the GetTweets endpoint */
type TweetsResponse struct {
	Tweets []*models.Tweet `json:"tweets"`
	Total  int64           `json:"total"`
}
