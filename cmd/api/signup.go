package main

import (
	"errors"
	"net/http"
	"shareapp/internal/data"
	"shareapp/internal/domain"
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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &domain.User{
		PublicID: publicID,
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if domain.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	dbUser, err := app.queries.CreateUser(r.Context(), data.CreateUserParams{
		PublicID:     user.PublicID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash(),
	})

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			app.serverErrorResponse(w, r, errors.New("duplicate email"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": dbUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
