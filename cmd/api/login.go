package main

import (
	"encoding/json"
	"net/http"
	"shareapp/internal/data"
)

func (app *application) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.queries.GetUserByEmailAuth(r.Context(), input.Email)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err.Error())
		return
	}

	domainUser := data.MapUserDBToDomain(&user)

	match, err := domainUser.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.errorResponse(w, r, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := app.jwtMaker.CreateToken(user.Email, user.ID)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]string{
		"username": user.Username,
		"email":    user.Email,
		"token":    token,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/v1/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
