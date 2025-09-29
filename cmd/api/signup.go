package main

import (
	"encoding/json"
	"net/http"
	"shareapp/internal/data"
	"shareapp/utils"
)

func (app *application) handleSignupPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		user, err := app.queries.CreateUser(r.Context(), data.CreateUserParams{
			Username:     input.Username,
			Email:        input.Email,
			PasswordHash: hashedPassword,
		})

		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		data := map[string]string{
			"username": user.Username,
			"email":    user.Email,
		}

		err = app.writeJSON(w, http.StatusCreated, data, nil)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

	}
}
