package relational

import (
	"time"
)

/* User model for the postgreSQL DB */
type Tweet struct {
	Id      uint64    `gorm:"primarykey"`
	Message string    `gorm:"not null"`
	Date    time.Time `gorm:"not null"`
	Active  bool      `gorm:"not null;default:true"`
	UserId  uint64    `gorm:"not null"`
}
