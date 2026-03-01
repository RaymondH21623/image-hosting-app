package main

import "net/http"

func (app *application) handleMeGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected route"))
}
