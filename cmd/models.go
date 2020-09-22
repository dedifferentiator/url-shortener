package main

//Url instance of url record in db
type Url struct {
	ID       int
	ShortUrl string
	OrigUrl  string
}

//Create creates an url instance in db
func (u Url) Create(db *DB) (string, error) {
	url, err := db.InsertUrl(u)
	if err != nil {
		return "", err
	}

	return url, nil
}

//GetOrigUrl returns OrigUrl by its shortUrl
func GetOrigUrl(shortUrl string, db *DB) (string, error) {
	url, err := db.GetOrigUrl(shortUrl)
	if err != nil {
		return "", err
	}

	return url, nil
}

//idToShortUrl transformes id into url using base
func idToShortUrl(id int) string {
	digits := []int{}

	for c := id; c > 0; {
		mod := c % base
		c = c / base
		digits = append(digits, mod)
	}

	url := []byte{}
	for _, digit := range digits {
		url = append(url, alphabet[digit])
	}

	return string(url)
}
