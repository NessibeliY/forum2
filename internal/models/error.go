package models

type ErrorPage struct {
	TextError  string
	StatusCode int
}

type ErrorResponse struct {
	Error string
}
