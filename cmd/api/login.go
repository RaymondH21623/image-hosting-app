package main

import (
	"encoding/json"
	"net/http"
	"shareapp/utils"
)

func (app *application) handleLoginPost() http.HandlerFunc {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, err := app.queries.GetUserByEmailAuth(r.Context(), input.Email)
		if err != nil {
			http.Error(w, "user not found", http.StatusInternalServerError)
			return
		}
		if utils.CheckPassword(input.Password, user.PasswordHash) != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := app.jwtMaker.CreateToken(input.Email)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"username": user.Username,
			"email":    user.Email,
			"token":    token,
		}

		//user := data.User{}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			HttpOnly: true,
		})
		json.NewEncoder(w).Encode(resp)
	}
}
