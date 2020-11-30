package models

import (
	"fmt"
	_ "fmt"

	null "gopkg.in/volatiletech/null.v6"
)

// Base model properties are shared with all models
type Base struct {
	ID        int       `json:"id"`
	CreatedAt null.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt"`
}

// User represents user in lowkey system
type User struct {
	Base

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Website   string `json:"website"`
}

// Users represents collection of user
type Users struct {
	Users []User `json:"users"`
}

// Name returns Full name of user by concating first & last name
func (u User) Name() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
