package main

import (
	"net/http"
)

// App stores gloabls & context vars
type App struct{}

// Handler executes middleware and provides a ref to our
// global App struct
type Handler func(a *App, w http.ResponseWriter, r *http.Request) error

// Handle executes all our middleware. Each middleware
// function accepts an extra argument: the `a *app.App`,
// which is our global context variable
func (a *App) Handle(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			err := handler(a, w, r)
			if err != nil {
				// Lazily handle errors for now
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}

// HomeHandler renders the home page
func HomeHandler(a *App, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Hello"))
	return nil
}
