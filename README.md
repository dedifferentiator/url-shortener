![Tests](https://github.com/dedifferentiator/url-shortener/workflows/Tests/badge.svg)
[![Maintainability](https://api.codeclimate.com/v1/badges/5d01cf982cc60af85b4e/maintainability)](https://codeclimate.com/github/dedifferentiator/url-shortener/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedifferentiator/url-shortener)](https://goreportcard.com/report/github.com/dedifferentiator/url-shortener)

# url-shortner

## Implemented solutions
  - [x] The requests to shortened URLs should be redirected to their original URL (status 302)
    or return 404 for unknown URLs.
  - [x] Simple HTML form should be served on the index page where users can input URL and
    retrieve the shortened version from server.
  - [X] All of the implemented HTTP handlers should have unit tests.
  - [ ] (optional) All shortened URLs should be persisted locally to a file using simple
    storage methods (SQLite, BoltDB, CSV..).
  - [ ] (optional) The redirect requests should be cached in memory for certain amount of time.

## Installation
  - `git clone https://github.com/dedifferentiator/url-shortener.git`
  - `cd url-shortener`
  - `docker-compose up --build`
  
## Configuration
  There are 2 configs for docker-compose `.env_serv` and `.env_postgres`. Their samples are available in the repository.
  ### .env_serv
  - `SERV_DRIVER_DB` - database driver
  - `SERV_CONN_DB` - connection string 
  - `SERV_PORT` - port which server will listen to
  - `SERV_DOMAIN` - _[optional]_ domain name of the host on which the url-shortener will run on. Used for generating shortened links. Default value `localhost`
  ## .env_postgres
  Postgres envars usage explained here https://hub.docker.com/_/postgres

## Link shortening
  Firstly, we want to decide what symbols we want to use for generating short-links. For this setup was choosen [`[a-Z0-9]`](https://github.com/dedifferentiator/url-shortener/blob/a2d3d16215c6d798a098b7e7144bcf3ab9a0086d/cmd/init.go#L13).
  
  Next we [request](https://github.com/dedifferentiator/url-shortener/blob/a2d3d16215c6d798a098b7e7144bcf3ab9a0086d/cmd/db.go#L124) `id SERIAL` from DB and after retrieving it, generate short-link using id value and insert a record with id, shortened and original links into the [table](https://github.com/dedifferentiator/url-shortener/blob/a2d3d16215c6d798a098b7e7144bcf3ab9a0086d/cmd/db.go#L8).
  
  Shortening links is quite simple, having serial id - integer(implicitly unsigned), we convert it from base10 to base62 (length of the pull of symbols) and then translate it to the shortened link, using the converted number digits as indices of the symbols pull.
  
  According to the above, there's even no need in `short_url` field in the table, because it's easy to convert from the shortened url to base62, but it was decided to store the short link anyway, in case of switching base to another number, different from 62. Although afterwhile I realised that it shouldn't be difficult to start storing shortened links in case we'll want to move to another base, but before it there's no actual need. But it was too late and too close to deadline :)

## Possible optimizations
 - forbid some links which contain bad words in them
 - refill holes in ID sequence in the db with the new urls
 - link-to-shorten size is not limited, but probably it should be
 - service tries to redirect even to invalid links, but probably it's worth to add a link validator
