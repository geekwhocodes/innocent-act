package main

import (
	"net/http"
	"regexp"

	"github.com/geekwhocodes/innocent-relay/models"
	"github.com/labstack/echo"
)

var (
	strInputSize = 128
	emailRegexp  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// RegisterHandlers d
func registerHTTPHandlers(e *echo.Echo) {
	//home page
	e.GET("/", handleIndexPage)
	e.GET("health", handlerHealth)
	// users
	e.POST("api/users", handlerCreateUser)
	e.GET("api/users/:id", handlerGetUser)
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
	return c.JSON(http.StatusOK, models.OkResponse{Data: "ok"})
}
