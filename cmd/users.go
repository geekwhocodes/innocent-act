package main

import (
	"errors"
	"fmt"
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
	// TODO : validate user object
	if err := validateUserFields(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := app.store.CreateUser(u)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, models.OkResponse{Data: user})
}

func handlerGetUser(c echo.Context) error {
	app := c.Get("app").(*App)
	userID := c.Param("id")
	if len(userID) <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user id is not valid")
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error while parsing user id: %v", err))
	}
	user, err := app.store.GetUser(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, models.OkResponse{Data: user})
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
	return c.JSON(http.StatusOK, models.OkResponse{Data: users})
}

// validate validates user fields by checking every field against business/common rules
func validateUserFields(u *models.User) error {
	if len(u.Name) == 0 || len(u.Name) > strInputSize {
		return errors.New("invalid or empty user name")
	}
	if len(u.Email) == 0 || len(u.Email) > strInputSize ||
		!emailRegexp.MatchString(u.Email) {
		return errors.New("invalid or empty user email")
	}
	if err := ValidatEmaileHost(u.Email); err != nil {
		return errors.New("invalid user email")
	}
	return nil
}
