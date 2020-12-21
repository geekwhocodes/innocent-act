package sqlite

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/geekwhocodes/innocent-relay/models"
	"github.com/labstack/echo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Options sqlite db config
type Options struct {
	Name string `koanf:"name"`
	Path string `koanf:"path"`
}

// Client implements `store.Store`
type Client struct {
	db *gorm.DB
}

// GetEmail .s
func (store *Client) GetEmail() string {
	return "Email created"
}

// CreateUser .s
func (store *Client) CreateUser(u *models.User) (*models.User, error) {
	user := models.User{}
	if result := store.db.Where("email = ?", u.Email).First(&user); result.Error != nil {
		fmt.Print(result.Error)
	}
	if user.ID != 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}
	fmt.Println("Creating user.", u.Name)
	if result := store.db.Create(u); result.Error != nil {
		fmt.Print(result.Error)
		return nil, result.Error
	}
	return u, nil
}

// GetUser retuens user by id
func (store *Client) GetUser(userID int) (*models.User, error) {
	user := models.User{}
	if result := store.db.Where("id = ?", userID).First(&user); result.Error != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return &user, nil
}

// GetAllUsers returns all users in db
func (store *Client) GetAllUsers() ([]*models.User, error) {
	users := []*models.User{}
	result := store.db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
		return users, result.Error
	}
	return users, nil
}

// GetPaginatedUsers returns paginated results
func (store *Client) GetPaginatedUsers(page int) ([]*models.User, error) {
	limit := 10
	offset := limit * (page - 1)

	//sqlStatement := `SELECT * FROM "users" ORDER BY "id" LIMIT $1 OFFSET $2`
	//rows, err := store.db.Queryx(sqlStatement, limit, offset)
	users := []*models.User{}
	result := store.db.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
		return users, result.Error
	}
	return users, nil
}

// Close .
func (store *Client) Close() error {
	sqlDB, err := store.db.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// NewDb .
func NewDb(options *Options) (*Client, error) {
	dbPath := filepath.Join(options.Path, options.Name+".db")
	fmt.Println("initializing db at ", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	return &Client{
		db: db,
	}, nil
}
