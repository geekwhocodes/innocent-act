package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/labstack/echo"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// App represents lowkey app's context
type App struct {
	//log *log.Logger
	db *sql.DB
	// channel to handle gracefull shutdown
	signalChan chan os.Signal
}

var (
	db      *sql.DB
	build   string
	version string
)

func init() {
	// TODO : check schema
}

func main() {
	fmt.Println("This is my lowkey app.")
	app := &App{
		db: initDb(),
	}
	sigs := make(chan os.Signal, 1)
	quitSing := make(chan bool, 1)
	// Start HTTP server
	server := initHTTPServer(app)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Quit gracefully
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		fmt.Println("Shutting down web server")
		server.Shutdown(ctx)
		defer cancel()
		fmt.Print("Closing Db")
		app.db.Close()
		fmt.Print("Exiting...")
		quitSing <- true
	}()
	<-quitSing
}
