package response

import mr "models/result"

/* UserTweetsResponse Response model for the GetUsersTweets endpoint */
type UserTweetsResponse struct {
	Tweets []*mr.UserTweet `json:"tweets"`
	Total  int64           `json:"total"`
}
