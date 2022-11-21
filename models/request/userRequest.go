package request

import "time"

/* User request */
type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	BirthDate time.Time `json:"birthDate,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Banner    string    `json:"banner,omitempty"`
	Biography string    `json:"biography,omitempty"`
	Location  string    `json:"location,omitempty"`
	WebSite   string    `json:"webSite,omitempty"`
}
