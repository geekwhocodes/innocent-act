package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/geekwhocodes/innocent-relay/internal/store"
	"github.com/geekwhocodes/innocent-relay/internal/store/providers/sqlite"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/stuffbin"
	flag "github.com/spf13/pflag"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func initFlags() {
	//
	// Use the POSIX compliant pflag lib instead of Go's flag lib.
	f := flag.NewFlagSet("config", flag.ExitOnError)
	f.Usage = func() {
		// show help
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// register commands
	f.Bool("version", false, "current version of the build")
	f.Bool("sample-config", false, "generate new config file.")
	f.String("config-file", "config.yaml", "Configuration file path. Generate new one by running --sample-config")
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error reading flags: %v", err)
	}

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

}

func initFS(dir string) stuffbin.FileSystem {
	// Get the executable's path.
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %v", err)
	}

	// Load the static files stuffed in the binary.
	fs, err := stuffbin.UnStuff(path)
	if err != nil {
		// Running in local mode. Load local assets into
		// the in-memory stuffbin.FileSystem.
		log.Printf("unable to initialize embedded filesystem: %v", err)
		files := []string{
			"config.sample.yaml",
			"web/dist/web:web",
			"web/dist/favicon.ico:/web/favicon.ico",
		}

		fs, err = stuffbin.NewLocalFS("/", files...)
		if err != nil {
			log.Fatalf("failed to initialize local fs: %v", err)
		}
	}
	return fs
}

func generateNewConfig() error {
	if _, err := os.Stat("config.toml"); !os.IsNotExist(err) {
		return errors.New("config.yaml already exists. Remove it to generate a new one")
	}

	// Initialize the static file system into which all
	// required static assets (.sql, .js files etc.) are loaded.
	fs := initFS("")
	b, err := fs.Read("config.sample.yaml")
	if err != nil {
		return fmt.Errorf("error reading sample config: %v", err)
	}

	return ioutil.WriteFile("config.yaml", b, 0644)
}

func initConfig(k *koanf.Koanf) {
	confFile := k.String("config-file")
	yamlFile := file.Provider(confFile)
	if err := k.Load(yamlFile, yaml.Parser()); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Config file not found. You can generate new one by running --sample-config command")
		}
		log.Fatalf("error loading file: %v", err)
	}
}

func initDbStore() store.Store {
	switch provider := k.String("db.provider"); provider {
	case "sqlite":
		{
			// Init sqlite with sqlite options
			var options sqlite.Options
			k.Unmarshal("db.sqlite", &options)
			dbStore, err := sqlite.NewDb(&options)
			if err != nil {
				log.Fatal(err)
			}
			return dbStore
		}
	default:
		log.Fatalf("unknown db provider. select sqlite")
	}
	return nil
}

func initHTTPServer(app *App) *echo.Echo {

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

	// Inject our App context to all http handlers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	// Register file server to serve static content
	fileServer := app.fs.FileServer()
	//e.GET("/public/*", echo.WrapHandler(fileServer))
	e.GET("/web/*", echo.WrapHandler(fileServer))

	registerHTTPHandlers(e)

	// Start the server.
	go func() {
		if err := e.Start(k.String("host") + ":" + k.String("port")); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println("HTTP server shut down")
			} else {
				log.Fatalf("error starting HTTP server: %v", err)
			}
		}
	}()

	return e
}

//
func quit(app *App) {
	sig := <-app.signalChan
	fmt.Println("Quiting due to :", sig)
	fmt.Println("Closing Db")
	app.store.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	fmt.Println("Shutting down web server")
	app.server.Shutdown(ctx)
	defer cancel()
	fmt.Println("Exiting...")
}

func getVersion(v string, b string) string {
	s := fmt.Sprintf(`
#      ____                                  __                __ 
#     /  _/___  ____  ____  ________  ____  / /_   ____ ______/ /_
#     / // __ \/ __ \/ __ \/ ___/ _ \/ __ \/ __/  / __ / ___/ __/
#   _/ // / / / / / / /_/ / /__/  __/ / / / /_   / /_/ / /__/ /_  
#  /___/_/ /_/_/ /_/\____/\___/\___/_/ /_/\__/   \__,_/\___/\__/
#  							  
#  %s %s`, v, b)
	return s
}
