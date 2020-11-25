package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/geekwhocodes/innocent-relay/models"
	"github.com/labstack/echo"
)

func handlerCreateUser(c echo.Context) error {
	app := c.Get("app").(*App)
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	user, err := createUser(*u, app.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error")
	}
	return c.JSON(http.StatusCreated, user)
}

func handlerGetUsers(c echo.Context) error {
	app := c.Get("app").(*App)
	users, err := getUsers(app.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error")
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser creates user in db
func createUser(u models.User, db *sql.DB) (models.User, error) {
	fmt.Println("Creating user.", u.Name())
	sqlStatement := "INSERT INTO users (firstname, lastname, email, website, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := db.Exec(sqlStatement, u.FirstName, u.LastName, u.Email, u.Website, time.Now().UTC().String(), time.Now().UTC().String())
	if err != nil {
		log.Fatal(err)
		return u, err
	}
	fmt.Println(*&res)
	return u, nil
}

// GetUsers returns all users
func getUsers(db *sql.DB) (models.Users, error) {
	sqlStatement := "SELECT * FROM users"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
		return models.Users{}, err
	}
	defer rows.Close()
	result := models.Users{}

	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Website, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		result.Users = append(result.Users, user)
	}
	return result, nil
}
