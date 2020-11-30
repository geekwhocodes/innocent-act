package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/geekwhocodes/innocent-relay/models"
	"github.com/labstack/echo"
)

func handlerCreateUser(c echo.Context) error {
	app := c.Get("app").(*App)
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	user, err := app.store.CreateUser(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error")
	}
	return c.JSON(http.StatusCreated, user)
}

func handlerGetUsers(c echo.Context) error {
	app := c.Get("app").(*App)
	pageStr := c.QueryParam("page")
	if pageStr == "" {
		users, err := app.store.GetAllUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error")
		}
		return c.JSON(http.StatusOK, users)
	}
	pageNo, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Fatal(err)
	}
	users, err := app.store.GetPaginatedUsers(pageNo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error")
	}
	// c.Response().Header().Add("X-Next-Link", c.Request().RequestURI)
	// c.Response().Header().Add("Access-Control-Expose-Headers", "X-Next-Link")
	return c.JSON(http.StatusOK, users)
}
