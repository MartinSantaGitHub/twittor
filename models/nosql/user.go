package nosql

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* User model for the mongo DB */
type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	LastName  string             `bson:"lastName"`
	BirthDate time.Time          `bson:"birthDate"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Avatar    string             `bson:"avatar"`
	Banner    string             `bson:"banner"`
	Biography string             `bson:"biography"`
	Location  string             `bson:"location"`
	WebSite   string             `bson:"webSite"`
}
