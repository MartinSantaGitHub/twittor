package response

import mr "models/request"

/* UserTweetsResponse is the response model for the GetFollowingTweets endpoint */
type UserTweetsResponse struct {
	Tweets []*mr.UserTweet `json:"tweets"`
	Total  int64           `json:"total"`
}
