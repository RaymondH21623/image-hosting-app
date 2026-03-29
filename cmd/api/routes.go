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
	router.Get("/v1/health", app.handleHealthGet)
	router.Get("/v1/", app.handleHelloGet)
	router.Post("/v1/signup", app.handleSignupPost)
	router.Post("/v1/login", app.handleCreateAuthenticationToken)
	router.Get("/v1/me", app.requireActivatedUser(app.handleMeGet))
	router.Post("/v1/media", app.requireActivatedUser(app.handleMediaPost))
	router.Get("/v1/media/{id}", app.requireActivatedUser(app.handleMediaGet))
	//router.Get("/v1/i/{id}", app.requireActivatedUserapp.serveMedia())
	router.Get("/v1/u/{id}", app.requireActivatedUser(app.handleMediaListGet))
	router.Post("/v1/tokens/activation", app.createActivationToken)
	router.Put("/v1/users/activated", app.handleActivateUserPut)

	return app.recoverPanic(app.authenticate(router))
}
