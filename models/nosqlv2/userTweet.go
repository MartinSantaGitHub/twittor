package nosqlv2

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* UserTweet Model of the User and Tweet results from the mongo DB */
type UserTweet struct {
	UserId          primitive.ObjectID `bson:"userId"`
	UserFollowingId primitive.ObjectID `bson:"userFollowingId"`
	Tweet           struct {
		Id      primitive.ObjectID `bson:"_id"`
		Message string             `bson:"message"`
		Date    time.Time          `bson:"date"`
	}
}
