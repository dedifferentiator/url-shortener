package main

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

type server struct {
	mux *echo.Echo
	db  *DB
	env envars
}

func main() {
	s := NewServer()

	t := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("static/new.html")),
	}
	s.mux.Renderer = t

	s.mux.File("/", "static/index.html")
	//TODO create page for /new
	s.mux.POST("/new", s.NewHandler)
	s.mux.GET("/*", s.RedirectHandler)

	s.mux.Logger.Fatal(s.mux.Start(":" + s.env.servPort))
}
