package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/geekwhocodes/innocent-relay/internal/store"
	"github.com/knadh/koanf"
	echo "github.com/labstack/echo"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// App represents lowkey app's context
type App struct {
	//log *log.Logger
	server *echo.Echo
	store  store.Store
	// channel to handle gracefull shutdown
	signalChan chan os.Signal
	quitSign   chan bool
}

var (
	//db      *sql.DB
	build      string
	version    string
	configPath string
	k          = koanf.New(".")
)

func init() {
	initFlags()
	initConfig()
}

func main() {
	fmt.Println("This is my lowkey app.")

	app := &App{
		store: initDbStore(),
	}

	// Start HTTP server
	app.server = initHTTPServer(app)

	app.signalChan = make(chan os.Signal, 1)
	app.quitSign = make(chan bool, 1)
	signal.Notify(app.signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Quit gracefully
	quit(app)
}
