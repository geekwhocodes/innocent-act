package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/geekwhocodes/innocent-relay/internal/store"
	"github.com/geekwhocodes/innocent-relay/internal/store/providers/sqlite"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	flag "github.com/spf13/pflag"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func initFlags() {
	//
	// Use the POSIX compliant pflag lib instead of Go's flag lib.
	configCmd := flag.NewFlagSet("config", flag.ContinueOnError)
	configCmd.Usage = func() {
		fmt.Println(configCmd.FlagUsages())
		os.Exit(0)
	}
	// Path to one or more config files to load into koanf along with some config params.
	yamlFilePath := configCmd.String("f", "./config.yml", "path to one or more .yml config file")
	// Actually parse the flags
	if err := configCmd.Parse(os.Args[:1]); err != nil {
		log.Fatalf("error loading flags: %v", err)
	}
	yamlFile := file.Provider(*yamlFilePath)
	if err := k.Load(yamlFile, yaml.Parser()); err != nil {
		log.Fatalf("error loading file: %v", err)
	}
	// Watch the file and get a callback on change. The callback can do whatever,
	// like re-load the configuration.
	// File provider always returns a nil `event`.
	// yamlFile.Watch(func(event interface{}, err error) {
	// 	if err != nil {
	// 		log.Printf("watch error: %v", err)
	// 		return
	// 	}

	// 	log.Println("config changed. Reloading ...")
	// 	k.Load(yamlFile, yaml.Parser())
	// 	k.Print()
	// })

	// // Block forever (and manually make a change to mock/mock.json) to
	// // reload the config.
	// log.Printf("listening changes to %s to live reload", *yamlFilePath)
	// <-make(chan bool)
}

func initConfig() *Config {
	// Create config structure
	config := &Config{}
	if err := k.Unmarshal("host", &config.Host); err != nil {
		log.Fatalf("error loading app config: %v", err)
	}
	if err := k.Unmarshal("port", &config.Port); err != nil {
		log.Fatalf("error loading app config: %v", err)
	}
	if err := k.Unmarshal("db", &config.Db); err != nil {
		log.Fatalf("error loading app config: %v", err)
	}

	return config
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

//
func quit(app *App) {
	sig := <-app.signalChan
	fmt.Println("Quiting due to :", sig)
	fmt.Println("Closing Db")
	//app.store.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	fmt.Println("Shutting down web server")
	app.server.Shutdown(ctx)
	defer cancel()
	fmt.Println("Exiting...")
}
