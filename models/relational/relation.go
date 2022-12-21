package relational

/* Relation Model for saving a relation between an user with another */
type Relation struct {
	Active bool `gorm:"not null;default:true"`
}
