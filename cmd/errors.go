package main

type DBError interface {
	Error() string
}

type ShortUrlIsEmptyErr struct{}

var _ DBError = (*ShortUrlIsEmptyErr)(nil)

func (e ShortUrlIsEmptyErr) Error() string {
	return "ShortUrl has to be a non-empty string"
}

type OrigUrlIsEmptyErr struct{}

var _ DBError = (*OrigUrlIsEmptyErr)(nil)

func (e OrigUrlIsEmptyErr) Error() string {
	return "OrigUrl has to be a non-empty string"
}

type CannotRetrieveIDErr struct{}

var _ DBError = (*CannotRetrieveIDErr)(nil)

func (e CannotRetrieveIDErr) Error() string {
	return "Cannot retrieve new id for inserting message"
}

type CannotRetrieveOrigUrlErr struct{}

var _ DBError = (*CannotRetrieveOrigUrlErr)(nil)

func (e CannotRetrieveOrigUrlErr) Error() string {
	return "Cannot retrieve orig_url for inserting message"
}
