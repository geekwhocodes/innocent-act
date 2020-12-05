package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/geekwhocodes/innocent-relay/internal/store"
	"github.com/knadh/koanf"
	"github.com/knadh/stuffbin"
	echo "github.com/labstack/echo"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// App represents lowkey app's context
type App struct {
	//log *log.Logger
	server *echo.Echo
	store  store.Store
	fs     stuffbin.FileSystem
	// channel to handle gracefull shutdown
	signalChan chan os.Signal
	quitSign   chan bool
}

var (
	k = koanf.New(".")

	buildString   string
	versionString string
)

func init() {
	initFlags()
}

func main() {
	// Handle arg flags/commands
	if k.Bool("version") {
		fmt.Print(getVersion(versionString, buildString))
		fmt.Println()
		fmt.Println()
		os.Exit(0)
	}

	// Handle arg flags/commands
	if k.String("config-file") == "" {
		fmt.Print(getVersion(versionString, buildString))
		fmt.Println()
		fmt.Println()
		os.Exit(0)
	}

	// Generate new config.
	if k.Bool("sample-config") {
		if err := generateNewConfig(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("new config 'config.yaml' generated")
		os.Exit(0)
	}

	initConfig(k)

	app := &App{
		store: initDbStore(),
		fs:    initFS(""),
	}

	// Start HTTP server
	app.server = initHTTPServer(app)

	app.signalChan = make(chan os.Signal, 1)
	app.quitSign = make(chan bool, 1)
	signal.Notify(app.signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Quit gracefully
	quit(app)
}
