package store

import (
	"github.com/geekwhocodes/innocent-relay/models"
)

// Store is our application data store interface
type Store interface {
	CreateUser(u *models.Users) (*models.Users, error)
	GetAllUsers() ([]*models.Users, error)
	GetPaginatedUsers(page int) ([]*models.Users, error)
	GetEmail() string
	Close() error
}
