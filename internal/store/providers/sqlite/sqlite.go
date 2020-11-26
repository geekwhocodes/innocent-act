package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gchaincl/dotsql"
	"github.com/geekwhocodes/innocent-relay/models"
)

// Options .
type Options struct {
	Name string `koanf:"name"`
	Path string `koanf:"path"`
}

// Client implements `store.Store`
type Client struct {
	db *sql.DB
}

// GetEmail .s
func (store *Client) GetEmail() string {
	return "Email created"
}

// CreateUser .s
func (store *Client) CreateUser(u *models.User) (*models.User, error) {
	fmt.Println("Creating user.", u.Name())
	sqlStatement := "INSERT INTO users (firstname, lastname, email, website, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := store.db.Exec(sqlStatement, u.FirstName, u.LastName, u.Email, u.Website, time.Now().UTC().String(), time.Now().UTC().String())
	if err != nil {
		log.Fatal(err)
		return u, err
	}
	return u, nil
}

// GetAllUsers creates user in db
func (store *Client) GetAllUsers() ([]*models.User, error) {
	sqlStatement := "SELECT * FROM users"
	rows, err := store.db.Query(sqlStatement)
	result := []*models.User{}
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Website, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		result = append(result, user)
	}
	return result, nil
}

// Close .
func (store *Client) Close() error {
	if err := store.db.Close(); err != nil {
		return err
	}
	return nil
}

// NewDb .
func NewDb(options *Options) (*Client, error) {
	dbPath := filepath.Join(options.Path, options.Name+".db")
	fmt.Println("initializing db at ", dbPath)
	// remove existing db file
	if err := validatePath(dbPath); err == nil {
		if err = os.Remove(dbPath); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err := createSchema(db); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

func validatePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func createSchema(db *sql.DB) error {
	dot, err := dotsql.LoadFromFile("schema.sql")
	if err != nil {
		log.Fatal(err)
		return err
	}
	// check table exists
	stmt, err := dot.Prepare(db, "create-users-table")
	if err != nil {
		log.Fatal(err)
		return err
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	if _, err := dot.Exec(db, "create-user", "Ganesh", "raskar", "email@mail.com", "https://kl.in", nil, nil); err != nil {
		return err
	}

	return nil
}
