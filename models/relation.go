package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/* Relation Model for saving a relation between an user with another */
type Relation struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserId         primitive.ObjectID `bson:"userId" json:"userId"`
	UserRelationId primitive.ObjectID `bson:"userRelationId" json:"userRelationId"`
	Active         bool               `bson:"active" json:"-"`
}
