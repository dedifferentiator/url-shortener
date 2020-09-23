package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var s = NewServer()
var shortUrls = []string{}

func TestMain(m *testing.M) {
	// creating a couple of records in db for further testing
	urls := []Url{
		{OrigUrl: "http://github.com/"},
		{OrigUrl: "https://docs.docker.com/compose/compose-file/"},
		{OrigUrl: "https://github.com/dedifferentiator/url-shortner"},
	}

	for _, url := range urls {
		sUrl, err := s.db.InsertUrl(url)
		if err != nil {
			log.Fatalln(err)
		}

		shortUrls = append(shortUrls, sUrl)
	}

	os.Exit(m.Run())
}

func TestNewHandler(t *testing.T) {
	// Setup
	formTrue := []url.Values{
		{"url": {"http://github.com/"}},
	}

	formFalse := []url.Values{
		{"url": {}},
		{},
	}

	// tests for valid form
	for _, form := range formTrue {
		// create request
		req := httptest.NewRequest(http.MethodPost, "/new", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

		rec := httptest.NewRecorder()
		c := s.mux.NewContext(req, rec)

		// assertions
		if assert.NoError(t, s.NewHandler(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	}

	// tests for non-valid form
	for _, form := range formFalse {
		// create request
		req := httptest.NewRequest(http.MethodPost, "/new", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

		rec := httptest.NewRecorder()
		c := s.mux.NewContext(req, rec)

		// assertions
		if assert.NoError(t, s.NewHandler(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "Something went wrong... Please, try again later!", rec.Body.String())
		}
	}
}

func TestRedirectHandler(t *testing.T) {

	urlsTrue := shortUrls
	urlsFalse := []string{
		"__test", "-", "-1", "Ы",
	}

	// tests for valid short urls
	for _, url := range urlsTrue {
		// create request
		req := httptest.NewRequest(http.MethodGet, "/"+url, nil)
		rec := httptest.NewRecorder()
		c := s.mux.NewContext(req, rec)

		// assertions
		if assert.NoError(t, s.RedirectHandler(c)) {
			assert.Equal(t, http.StatusFound, rec.Code)
		}
	}

	// tests for non-valid short urls
	for _, url := range urlsFalse {
		// create request
		req := httptest.NewRequest(http.MethodGet, "/"+url, nil)
		rec := httptest.NewRecorder()
		c := s.mux.NewContext(req, rec)

		// assertions
		if assert.NoError(t, s.RedirectHandler(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, "Oops...looks like this link is unvalid or doesn't exist! :(", rec.Body.String())
		}
	}

}
