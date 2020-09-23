package main

import (
	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE IF NOT EXISTS urls (
	id SERIAL,
	short_url VARCHAR(32) PRIMARY KEY,
	orig_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);
`

//DB database connector
type DB struct {
	DB     *sqlx.DB
	Driver string
	Conn   string
}

//OpenConn opens connection to database
func (d *DB) OpenConn() (*sqlx.DB, error) {
	db, err := sqlx.Connect(d.Driver, d.Conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//Close closes connection to database
func (d *DB) Close() error {
	return d.DB.Close()
}

//AutoMigrate migrates table schema
func (d *DB) AutoMigrate() error {
	db, err := d.OpenConn()
	if err != nil {
		return err
	}
	d.DB = db
	defer d.Close()

	_, err = d.DB.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

//InsertReservedWords inserts reserved handlers' keywords into db
func (d *DB) InsertReservedWords(urls []Url) error {
	db, err := d.OpenConn()
	if err != nil {
		return err
	}
	d.DB = db
	defer d.Close()

	// maybe it's worth to create a separate field for marking records
	// with keywords
	for _, u := range urls {
		if u.ShortUrl == "" {
			return ShortUrlIsEmptyErr{}
		}

		tx := d.DB.MustBegin()

		_, err = tx.Exec(`INSERT INTO urls (
                            id, short_url, orig_url, created_at
                          )
                          VALUES (
                            nextval(
                              pg_get_serial_sequence('urls', 'id')
                            ), $1, $2, now()
                          ) ON CONFLICT DO NOTHING`,
			u.ShortUrl, u.ShortUrl)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

//InsertUrl inserts Url into db
func (d *DB) InsertUrl(u Url) (string, error) {
	db, err := d.OpenConn()
	if err != nil {
		return "", err
	}
	d.DB = db
	defer d.Close()

	if u.OrigUrl == "" {
		return "", OrigUrlIsEmptyErr{}
	}

	// it would be great to implement this transaction with pl/pgSQL
	// so the record will be added into db through one query
	tx := d.DB.MustBegin()

	// querying id value to calculate short_url based on it

	// SERIAL returns new value each time, it's being queried,
	// but this might be not really good idea to insert record
	// in 2-step transaction, having in mind situation
	// when following INSERT query fails constantly
	// for a really huge number of times and no records in db being made,
	// but the value of id will be unproportianally enourmous
	// and probably this might be abused in some way
	id := []int{}
	err = tx.Select(&id, `SELECT nextval(
                            pg_get_serial_sequence('urls', 'id')
                          ) AS id`)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	} else if len(id) != 1 {
		return "", CannotRetrieveIDErr{}
	}

	_, err = tx.Exec(`INSERT INTO urls (
                        id, short_url, orig_url, created_at
                      )
                      VALUES (
                        $1, $2, $3, now()
                      )`, id[0], idToShortUrl(id[0]), u.OrigUrl)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	return idToShortUrl(id[0]), tx.Commit()
}

//GetOrigUrl returns original url by its shortened url
func (d *DB) GetOrigUrl(shortUrl string) (string, error) {
	db, err := d.OpenConn()
	if err != nil {
		return "", err
	}
	d.DB = db
	defer d.Close()

	url := []string{}
	err = d.DB.Select(&url,
		"SELECT orig_url FROM urls WHERE short_url=$1",
		shortUrl)
	if err != nil {
		return "", err
	} else if len(url) != 1 {
		return "", CannotRetrieveOrigUrlErr{}
	}

	return url[0], nil
}
