package nosql

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* Tweet model for the mongo DB */
type Tweet struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	UserId  primitive.ObjectID `bson:"userId"`
	Message string             `bson:"message"`
	Date    time.Time          `bson:"date"`
	Active  bool               `bson:"active"`
}
