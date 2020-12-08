package models

import (
	"time"
)

//go:generate reform

// Users represents a row in users table.
//reform:users
type Users struct {
	ID        int64      `reform:"ID,pk"`
	FirstName string     `reform:"FirstName"`
	LastName  string     `reform:"LastName"`
	Email     string     `reform:"Email"`
	Website   string     `reform:"Website"`
	CreatedAt *time.Time `reform:"CreatedAt"`
	UpdatedAt *time.Time `reform:"UpdatedAt"`
}
