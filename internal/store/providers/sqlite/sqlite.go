package sqlite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/geekwhocodes/innocent-relay/models"
	"github.com/jmoiron/sqlx"
	"gopkg.in/volatiletech/null.v6"
)

// Options .
type Options struct {
	Name string `koanf:"name"`
	Path string `koanf:"path"`
}

// Client implements `store.Store`
type Client struct {
	db *sqlx.DB
}

// GetEmail .s
func (store *Client) GetEmail() string {
	return "Email created"
}

// CreateUser .s
func (store *Client) CreateUser(u *models.User) (*models.User, error) {
	fmt.Println("Creating user.", u.Name())
	sqlStatement, _ := store.db.Prepare("INSERT INTO users (firstname, lastname, email, website, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)")
	u.CreatedAt, u.UpdatedAt = null.NewTime(time.Now(), true), null.NewTime(time.Now(), true)
	result, err := sqlStatement.Exec(u.FirstName, u.LastName, u.Email, u.Website, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	return u, nil
}

// GetAllUsers returns all users in db
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

// GetPaginatedUsers returns paginated results
func (store *Client) GetPaginatedUsers(page int) ([]*models.User, error) {
	limit := 10
	offset := limit * (page - 1)

	sqlStatement := `SELECT * FROM "users" ORDER BY "id" LIMIT $1 OFFSET $2`
	rows, err := store.db.Queryx(sqlStatement, limit, offset)
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
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	if err := createSchema(db); err != nil {
		log.Fatalln(err)
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

func createSchema(db *sqlx.DB) error {
	var schema = `
		DROP TABLE IF EXISTS user;
		CREATE TABLE users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				firstname VARCHAR(64) NULL,
				lastname VARCHAR(64) NULL,
				email VARCHAR(64) NULL,
				website VARCHAR(32) NULL,
				createdAt DATE NULL,
				updatedAt DATE NULL
		);`
	db.MustExec(schema)
	return nil
}
