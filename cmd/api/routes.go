package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/v1/health", app.handleHealthGet())
	router.Get("/v1/", app.handleHelloGet())
	router.Post("/v1/signup", app.handleSignupPost())
	router.Post("/v1/login", app.handleLoginPost())
	router.Get("/v1/me", app.authMiddleware(app.handleMeGet()))
	router.Post("/v1/media", app.authMiddleware(app.handleMediaPost()))
	// router.Get("/v1/media/{id}", app.authMiddleware(app.handleMediaGet()))
	// router.Get("/v1/i/{id}", app.serveMedia())

	return router
}
