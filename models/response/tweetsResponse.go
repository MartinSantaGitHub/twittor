package response

import mr "models/request"

/* TweetsResponse is the response model for the GetTweets and GetFollowingTweets endpoints */
type TweetsResponse struct {
	Tweets []*mr.Tweet `json:"tweets"`
	Total  int64       `json:"total"`
}
