package main

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

//TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

//Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//NewHandler handles request for creating new url
func (s *server) NewHandler(c echo.Context) error {
	p, err := c.FormParams()
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Something went wrong... Please, try again later!")
	}

	u := Url{
		OrigUrl: p.Get("url"),
	}

	shortUrl, err := u.Create(s.db)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Something went wrong... Please, try again later!")
	}

	link := s.env.domain + "/" + shortUrl
	return c.Render(http.StatusCreated, "new", link)
}

//RedirectHandler handles request with short urls
func (s *server) RedirectHandler(c echo.Context) error {
	shortUrl := strings.TrimLeft(c.Request().URL.Path, "/")

	url, err := GetOrigUrl(shortUrl, s.db)
	if err != nil {
		return c.String(http.StatusNotFound,
			"Oops...looks like this link is unvalid or doesn't exist! :(")
	}

	return c.Redirect(http.StatusFound, url)
}
