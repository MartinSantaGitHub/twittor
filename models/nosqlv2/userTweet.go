package nosqlv2

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* UserTweet Model of the User and Tweet results from the mongo DB */
type UserTweet struct {
	Id             primitive.ObjectID `bson:"_id"`
	UserId         primitive.ObjectID `bson:"userId"`
	UserRelationId primitive.ObjectID `bson:"userRelationId"`
	Tweet          struct {
		Id      primitive.ObjectID `bson:"_id"`
		Message string             `bson:"message"`
		Date    time.Time          `bson:"date"`
	}
}
