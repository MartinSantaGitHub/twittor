package nosqlv2

import "go.mongodb.org/mongo-driver/bson/primitive"

/* Relation Model for saving a relation between an user with another */
type Relation struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	UserId         primitive.ObjectID `bson:"userId"`
	UserRelationId primitive.ObjectID `bson:"userRelationId"`
	Active         bool               `bson:"active"`
}
