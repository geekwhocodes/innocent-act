package store

import (
	"github.com/geekwhocodes/innocent-relay/models"
)

// Store is our application data store interface
type Store interface {
	CreateUser(u *models.User) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetEmail() string
	Close() error
}
