package request

import "time"

/* UserTweet request */
type UserTweet struct {
	Id             string `json:"id,omitempty"`
	UserId         string `json:"userId,omitempty"`
	UserRelationId string `json:"userRelationId,omitempty"`
	Tweet          struct {
		Id      string    `json:"id,omitempty"`
		Message string    `json:"message,omitempty"`
		Date    time.Time `json:"date,omitempty"`
	}
}
