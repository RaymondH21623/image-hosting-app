package server

import (
	"encoding/json"
	"net/http"
	"shareapp/internal/db"
	"shareapp/utils"
)

func (s *Server) handleHealthGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func (s *Server) handleHelloGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}
}

func (s *Server) handleSignupPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u UserReq

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			http.Error(w, "failed to hash password", http.StatusInternalServerError)
			return
		}

		user, err := s.queries.CreateUser(r.Context(), db.CreateUserParams{
			Username:     u.Username,
			Email:        u.Email,
			PasswordHash: hashedPassword,
		})

		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		resp := UserRes{
			Username: user.Username,
			Email:    user.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) handleLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u LoginReq

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, err := s.queries.GetUserByEmail(r.Context(), u.Email)
		if err != nil {
			http.Error(w, "user not found", http.StatusInternalServerError)
			return
		}
		if utils.CheckPassword(u.Password, user.PasswordHash) != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := s.jwtMaker.CreateToken(u.Email)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		resp := LoginRes{
			Token:    token,
			Username: user.Username,
		}
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

func (s *Server) authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := c.Value

	}
}
