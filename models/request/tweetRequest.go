package request

import "time"

/* Tweet request model */
type Tweet struct {
	Id      string    `json:"id,omitempty"`
	UserId  string    `json:"userId,omitempty"`
	Message string    `json:"message,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Active  bool      `json:"-"`
}
