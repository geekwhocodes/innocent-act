package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/geekwhocodes/innocent-relay/models"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/sqlite3"
)

// Options .
type Options struct {
	Name string `koanf:"name"`
	Path string `koanf:"path"`
}

// Client implements `store.Store`
type Client struct {
	db         *reform.DB
	connection *sql.DB
}

// GetEmail .s
func (store *Client) GetEmail() string {
	return "Email created"
}

// CreateUser .s
func (store *Client) CreateUser(u *models.Users) (*models.Users, error) {
	fmt.Println("Creating user.", u)
	sqlStatement, _ := store.connection.Prepare("INSERT INTO users (firstname, lastname, email, website, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)")
	now := time.Now()
	u.CreatedAt, u.UpdatedAt = &now, &now
	result, err := sqlStatement.Exec(u.FirstName, u.LastName, u.Email, u.Website, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	u.ID = int64(id)
	return u, nil
}

// GetAllUsers returns all users in db
func (store *Client) GetAllUsers() ([]*models.Users, error) {
	sqlStatement := "SELECT * FROM users"
	rows, err := store.db.Querier.Query(sqlStatement)
	result := []*models.Users{}
	if err != nil {
		//log.Fatal(err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.Users{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Website, &user.CreatedAt, &user.UpdatedAt); err != nil {
			//log.Fatal(err)
			log.Print(err)
		}
		result = append(result, user)
	}
	return result, nil
}

// GetPaginatedUsers returns paginated results
func (store *Client) GetPaginatedUsers(page int) ([]*models.Users, error) {
	limit := 10
	offset := limit * (page - 1)

	sqlStatement := `SELECT * FROM "users" ORDER BY "id" LIMIT $1 OFFSET $2`
	rows, err := store.db.Query(sqlStatement, limit, offset)
	result := []*models.Users{}
	if err != nil {
		log.Print(err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.Users{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Website, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Print(err)
		}
		result = append(result, user)
	}
	return result, nil
}

// Close .
func (store *Client) Close() error {
	if err := store.connection.Close(); err != nil {
		return err
	}
	return nil
}

// NewDb .
func NewDb(options *Options) (*Client, error) {
	dbPath := filepath.Join(options.Path, options.Name+".db")
	fmt.Println("initializing db at ", dbPath)
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Print(err)
	}
	//defer sqlDB.Close()
	// Use new *log.Logger for logging.
	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	// Create *reform.DB instance with simple logger.
	// Any Printf-like function (fmt.Printf, log.Printf, testing.T.Logf, etc) can be used with NewPrintfLogger.
	// Change dialect for other databases.
	db := reform.NewDB(sqlDB, sqlite3.Dialect, reform.NewPrintfLogger(logger.Printf))

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return &Client{
		db:         db,
		connection: sqlDB,
	}, nil
}
