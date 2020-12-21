package models

import (
	"strings"

	"gorm.io/gorm"
)

// OkResponse data interface while returning data to the client
type OkResponse struct {
	Data interface{} `json:"data"`
}

// User represents user in lowkey system
type User struct {
	gorm.Model

	Name    string `json:"name"`
	Email   string `json:"email" gorm:"index;unique;not null"`
	Website string `json:"website"`
}

// Users represents collection of user
type Users struct {
	Users []User `json:"users"`
}

// FirstName it splits name string by space char and
// returns first part of the split if it's lenght is greater than 3
// assuming that the first part is user's first name
func (u User) FirstName() string {
	for _, u := range strings.Split(u.Name, " ") {
		if len(u) > 3 {
			return u
		}
	}
	return u.Name
}
