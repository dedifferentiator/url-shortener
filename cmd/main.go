package main

import (
	"github.com/labstack/echo/v4"
)

type server struct {
	mux *echo.Echo
	db  *DB
	env envars
}

func main() {
	s := NewServer()

	s.mux.File("/", "static/index.html")
	s.mux.POST("/new", s.NewHandler)
	s.mux.GET("/*", s.RedirectHandler)

	s.mux.Logger.Fatal(s.mux.Start(":" + s.env.servPort))
}
