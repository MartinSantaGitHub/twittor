package relational

import (
	"time"
)

/* User model for the postgreSQL DB */
type User struct {
	Id        uint64    `gorm:"primarykey"`
	Name      string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	BirthDate time.Time `gorm:"not null"`
	Email     string    `gorm:"not null;uniqueIndex"`
	Password  string    `gorm:"not null"`
	Avatar    string
	Banner    string
	Biography string
	Location  string
	WebSite   string
	Tweets    []Tweet
	Following []User `gorm:"many2many:relations;"`
}
