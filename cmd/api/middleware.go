package main

import "net/http"

func (app *application) authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err == http.ErrNoCookie {
			http.Error(w, "no cookie", http.StatusUnauthorized)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := c.Value

		err = app.jwtMaker.VerifyToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
}
