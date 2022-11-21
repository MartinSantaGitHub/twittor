package request

/* Relation is the request model for saving a relation between an user with another */
type Relation struct {
	Id             string `json:"-"`
	UserId         string `json:"userId,omitempty"`
	UserRelationId string `json:"userRelationId,omitempty"`
	Active         bool   `json:"-"`
}
