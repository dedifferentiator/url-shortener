package main

import (
	"html/template"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

//pull of digits from which shortened urls will consist
const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base     = len(alphabet)
)

type envars struct {
	driverDB string
	connDB   string
	servPort string
	domain   string
}

//NewServer creates a new server instance
func NewServer() *server {
	s := server{}

	s.initEnv()
	s.mustInitDB()
	s.mux = echo.New()

	t := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("static/new.html")),
	}
	s.mux.Renderer = t

	return &s
}

//mustInitDB performs automigrations
func (s *server) mustInitDB() {
	s.db = &DB{
		Driver: s.env.driverDB,
		Conn:   s.env.connDB,
	}

	err := s.db.AutoMigrate()
	if err != nil {
		log.Fatalf("[init] ERROR: cannot auto-migrate database: %e\n", err)
	}

	urls := []Url{
		{
			ShortUrl: "new",
		},
	}

	err = s.db.InsertReservedWords(urls)
	if err != nil {
		log.Fatalf("[init] ERROR: cannot insert reserved keywords %e\n", err)
	}
}

//initEnv imports envars
func (s *server) initEnv() {
	s.env = envars{}

	s.env.driverDB = os.Getenv("SERV_DRIVER_DB")
	if s.env.driverDB == "" {
		log.Fatalln("[init] ERROR: SERV_DRIVER_DB is not set")
	}

	s.env.connDB = os.Getenv("SERV_CONN_DB")
	if s.env.connDB == "" {
		log.Fatalln("[init] ERROR: SERV_CONN_DB is not set")
	}

	s.env.servPort = os.Getenv("SERV_PORT")
	if s.env.servPort == "" {
		log.Fatalln("[init] ERROR: SERV_PORT is not set")
	}

	s.env.domain = os.Getenv("SERV_DOMAIN")
	if s.env.domain == "" {
		log.Println("[init] WARNING: SERV_DOMAIN is not set; using default value 'localhost'")
		s.env.domain = "localhost"
	}
	if s.env.domain == "localhost" {
		s.env.domain += ":" + s.env.servPort
	}
}
