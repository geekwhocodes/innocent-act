package main

import (
	"database/sql"
	"fmt"
	"os"

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
	initHTTPServer(app)
	// // Start HTTP server
	// server := initHTTPServer(app)

	// app.signalChan = make(chan os.Signal)
	// signal.Notify(app.signalChan, syscall.SIGINT, syscall.SIGQUIT)

	// // Quit gracefully
	// go func() {
	// 	quitSig := <-app.signalChan
	// 	fmt.Println(quitSig)
	// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// 	fmt.Println("Shutting down web server")
	// 	server.Shutdown(ctx)
	// 	defer cancel()
	// 	fmt.Print("Closing Db")
	// 	app.db.Close()
	// 	fmt.Print("Exiting...")
	// 	app.signalChan <- syscall.SIGINT
	// }()
	// <-app.signalChan
}
