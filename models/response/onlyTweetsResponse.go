package response

import "models"

/* OnlyTweetsResponse Response model for the GetUsersTweets endpoint */
type OnlyTweetsResponse struct {
	Tweets []*models.Tweet `json:"tweets"`
	Total  int64           `json:"total"`
}
