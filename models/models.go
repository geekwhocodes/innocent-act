package models

import "fmt"

// Name returns Full name of user by concating first & last name
func (u *Users) Name() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
