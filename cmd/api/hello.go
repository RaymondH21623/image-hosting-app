package main

import (
	"net/http"
)

func (app *application) handleHelloGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}
}
