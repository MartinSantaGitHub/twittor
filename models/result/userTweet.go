package result

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* UserTweet Model of the User and Tweet results  */
type UserTweet struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserId         primitive.ObjectID `bson:"userId" json:"userId,omitempty"`
	UserRelationId primitive.ObjectID `bson:"userRelationId" json:"userRelationId,omitempty"`
	Tweet          struct {
		Id      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
		Message string             `bson:"message" json:"message,omitempty"`
		Date    time.Time          `bson:"date" json:"date,omitempty"`
	}
}
