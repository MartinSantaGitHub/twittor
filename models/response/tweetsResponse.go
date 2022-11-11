package response

import "models"

type TweetsResponse struct {
	Tweets []*models.Tweet `json:"tweets"`
	Total  int64           `json:"total"`
}
