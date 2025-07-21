package core

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HandlerFunc func(c IContext)

type IMux interface {
	Use(middlewares ...func(http.Handler) http.Handler)

	Mount(pattern string, h http.Handler)
	Group(fn func(r IMux)) IMux
	Route(patter string, fn func(r IMux)) IMux

	Get(pattern string, handlerFn HandlerFunc)
	Post(pattern string, handlerFn HandlerFunc)
	Put(pattern string, handlerFn HandlerFunc)
	Delete(pattern string, handlerFn HandlerFunc)
	ServeStaticFiles(url string, path string)

	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Mux struct {
	mux chi.Router
}

func NewMux() IMux {
	mux := chi.NewMux()
	return &Mux{mux: mux}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.mux.Use(middlewares...)
}

func (m *Mux) Get(pattern string, handlerFn HandlerFunc) {
	m.mux.Get(pattern, f(handlerFn))
}

func (m *Mux) Post(pattern string, handlerFn HandlerFunc) {
	m.mux.Post(pattern, f(handlerFn))
}

func (m *Mux) Put(pattern string, handlerFn HandlerFunc) {
	m.mux.Put(pattern, f(handlerFn))
}

func (m *Mux) Delete(pattern string, handlerFn HandlerFunc) {
	m.mux.Delete(pattern, f(handlerFn))
}

func (m *Mux) Mount(pattern string, h http.Handler) {
	m.mux.Mount(pattern, h)
}

func (m *Mux) Group(fn func(r IMux)) IMux {
	m.mux.Group(func(r chi.Router) {
		fn(&Mux{mux: r})
	})
	return m
}

func (m *Mux) Route(patter string, fn func(r IMux)) IMux {
	m.mux.Route(patter, func(r chi.Router) {
		fn(&Mux{mux: r})
	})
	return m
}

func (m *Mux) ServeStaticFiles(url string, path string) {
	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, path))
	m.Get(url, func(ctx IContext) {
		rctx := chi.RouteContext(ctx.Request().Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(ctx.Writer(), ctx.Request())
	})
}

func f(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		handler(ctx)
	}
}
