package core

import (
	"github.com/go-chi/render"
	"net/http"
)

// M is a convenience alias for quickly building a map structure that is going
// out to a responder. Just a short-hand.
type M map[string]any

// JSON marshals 'v' to JSON, and setting the Content-Type as application/json.
// Note that this does NOT auto-escape HTML.
func (c *context) JSON(status int, v any) {
	render.Status(c.Request(), status)
	render.JSON(c.Writer(), c.Request(), v)
}

// HTML writes a string to the response, setting the Content-Type as text/html.
func (c *context) HTML(status int, v string) {
	render.Status(c.Request(), status)
	render.HTML(c.Writer(), c.Request(), v)
}

// PlainText writes a string to the response, setting the Content-Type as
// text/plain.
func (c *context) PlainText(status int, v string) {
	render.Status(c.Request(), status)
	render.PlainText(c.Writer(), c.Request(), v)
}

// Data writes raw bytes to the response, setting the Content-Type as
// application/octet-stream.
func (c *context) Data(status int, v []byte) {
	render.Status(c.Request(), status)
	render.Data(c.Writer(), c.Request(), v)
}

// InternalServerError returns a HTTP 500 "Internal server error" response.
func (c *context) InternalServerError() {
	render.Status(c.Request(), http.StatusInternalServerError)
	render.JSON(c.Writer(), c.Request(), M{"message": "Something went wrong. If the issue persists, contact support."})
}

// BadRequest returns a HTTP 400 "Bad Request" response.
func (c *context) BadRequest(err error) {
	render.Status(c.Request(), http.StatusBadRequest)
	render.JSON(c.Writer(), c.Request(), M{"message": err.Error()})
}

// Forbidden returns a HTTP 403 "Forbidden" response.
func (c *context) Forbidden() {
	render.Status(c.Request(), http.StatusForbidden)
	render.JSON(c.Writer(), c.Request(), M{"message": "forbidden"})
}

// NotFound returns a HTTP 404 "Not Found" response.
func (c *context) NotFound() {
	render.Status(c.Request(), http.StatusNotFound)
	render.JSON(c.Writer(), c.Request(), M{"message": "not found"})
}

// Unauthorized returns a HTTP 401 "Unauthorized" response.
func (c *context) Unauthorized() {
	render.Status(c.Request(), http.StatusUnauthorized)
	render.JSON(c.Writer(), c.Request(), M{"message": "unauthorized"})
}

// NoContent returns a HTTP 204 "No Content" response.
func (c *context) NoContent() {
	render.NoContent(c.Writer(), c.Request())
}

// OK returns a HTTP 200 "OK" response.
func (c *context) OK(v any) {
	render.Status(c.Request(), http.StatusOK)
	render.JSON(c.Writer(), c.Request(), v)
}

// Created returns a HTTP 201 "Created" response.
func (c *context) Created(v any) {
	render.Status(c.Request(), http.StatusCreated)
	render.JSON(c.Writer(), c.Request(), v)
}
