package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Message string }{Message: "Aegis Start"})
	})
	e.Logger.Fatal(e.Start(":8000"))
}
