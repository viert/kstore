package store

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Error represents an error object for store-specific fields
type Error struct {
	StatusCode int
	Body       string
}

func (e Error) Error() string {
	return fmt.Sprintf("Unsuccessfull Response Code: %d\n\n%s", e.StatusCode, e.Body)
}

func responseError(resp *http.Response) Error {
	var body string

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("%v", err)
	} else {
		body = string(data)
	}

	return Error{
		StatusCode: resp.StatusCode,
		Body:       body,
	}
}

// IsFileNotFound checks if it's a StoreError and the response code is 404
func IsFileNotFound(err error) bool {
	se, ok := err.(Error)
	if !ok {
		return false
	}
	return se.StatusCode == 404
}

// IsForbidden checks if it's a StoreError and the response code is 403
func IsForbidden(err error) bool {
	se, ok := err.(Error)
	if !ok {
		return false
	}
	return se.StatusCode == 403
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
