package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

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

		claims, err := app.jwtMaker.VerifyToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		subject := claims.UserID

		userID, err := uuid.Parse(subject)

		if err != nil {
			http.Error(w, "invalid user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) recoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
