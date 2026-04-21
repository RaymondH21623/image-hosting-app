package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Use(middleware.Logger)
	router.Get("/v1/health", app.handleHealthcheck)
	router.Get("/v1/", app.handleHelloGet)
	router.Post("/v1/signup", app.handleSignupPost)
	router.Post("/v1/login", app.handleCreateAuthenticationToken)
	router.Get("/v1/me", app.requireActivatedUser(app.handleMeGet))
	router.Post("/v1/media", app.handleCreateMedia)
	router.Get("/v1/media/{id}", app.requireActivatedUser(app.handleShowMedia))
	//router.Get("/v1/i/{id}", app.requireActivatedUserapp.serveMedia())
	router.Get("/v1/{id}/media", app.requireActivatedUser(app.handleListUserMedia))
	router.Get("/v1/media", app.requireActivatedUser(app.handleListMedia))
	router.Post("/v1/tokens/activation", app.createActivationToken)
	router.Put("/v1/users/activated", app.handleActivateUserPut)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
