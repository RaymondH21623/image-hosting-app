package main

import (
	"net/http"
	"shareapp/internal/data"
	"shareapp/internal/validator"
	"shareapp/utils"
)

func (app *application) handleSignupPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	publicID, err := utils.GenerateID()
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := app.queries.CreateUser(r.Context(), data.CreateUserParams{
		PublicID:     publicID,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": dbUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
