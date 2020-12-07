package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// RegisterHandlers d
func registerHTTPHandlers(e *echo.Echo) {
	//home page
	e.GET("/", handleIndexPage)
	e.GET("health", handlerHealth)
	// users
	e.POST("api/users", handlerCreateUser)
	e.GET("api/users", handlerGetUsers)

	// static routes
	e.GET("/settings", handleIndexPage)
}

// handleIndex is the root handler that renders the Javascript frontend.
func handleIndexPage(c echo.Context) error {
	app := c.Get("app").(*App)

	b, err := app.fs.Read("/web/index.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Type", "text/html")
	return c.String(http.StatusOK, string(b))
}

func handlerHealth(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
