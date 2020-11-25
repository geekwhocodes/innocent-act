package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gchaincl/dotsql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func initDb() *sql.DB {
	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectDb() (*sql.DB, error) {
	var err error

	if fileExists("./lowkey.db") {
		if err := os.Remove("./lowkey.db"); err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", "./lowkey.db")
	if err := createSchema(db); err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}
	return db, nil
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

	if _, err := dot.Exec(db, "create-user", "Ganesh0", "raskar", "email@mail.com", "https://kl.kl", nil, nil); err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initHTTPServer(app *App) *echo.Echo {

	e := echo.New()
	e.Use(middleware.CORS())

	// Inject our App context to all http handlers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	registerHTTPHandlers(e)

	// Start the server.
	go func() {
		if err := e.Start(":8081"); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println("HTTP server shut down")
			} else {
				log.Fatalf("error starting HTTP server: %v", err)
			}
		}
	}()

	return e
}

func quit(c *chan os.Signal) {

}
