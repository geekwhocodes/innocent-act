package main

import (
	"github.com/labstack/echo"
)

// RegisterHandlers d
func registerHTTPHandlers(e *echo.Echo) {
	e.POST("api/users", handlerCreateUser)
	e.GET("api/users", handlerGetUsers)
}
