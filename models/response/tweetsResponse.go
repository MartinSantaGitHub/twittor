package response

import "models"

/* TweetsResponse Response model for the GetTweets and GetUsersTweets endpoints */
type TweetsResponse struct {
	Tweets []*models.Tweet `json:"tweets"`
	Total  int64           `json:"total"`
}
