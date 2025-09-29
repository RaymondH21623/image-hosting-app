package main

import "net/http"

func (app *application) handleMeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected route"))
	}
}
