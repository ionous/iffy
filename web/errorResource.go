package web

import (
	"io"
	"net/http"

	"golang.org/x/net/context"
)

// Error generates an error for every endpoint.
type Error struct {
	Message error
}

func (Error) Find(string) Resource {
	return nil
}

// Read a resource
func (e Error) Get(context.Context, http.ResponseWriter) error {
	return e.Message
}

// Alter a resource
func (e Error) Post(context.Context, io.Reader, http.ResponseWriter) error {
	return e.Message
}

// Receive a resource
func (e Error) Put(context.Context, io.Reader, http.ResponseWriter) error {
	return e.Message
}
