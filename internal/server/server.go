package server

import (
	"database/sql"
	"log"
	"net/http"

	"shareapp/internal/db"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	port    string
	queries *db.Queries
}

func New(port, dsn string) (*Server, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to db", "err", err)
	}
	defer conn.Close()

	s := &Server{
		router:  chi.NewRouter(),
		port:    port,
		queries: db.New(conn),
	}
	s.Routes()
	return s, nil
}

func (s *Server) Routes() {
	s.router.Get("/health", s.handleHealthGet())
	s.router.Get("/", s.handleHelloGet())
	s.router.Post("/signup", s.handleSignupPost())
	s.router.Post("/login", s.handleLoginPost())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

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
		w.Write([]byte("Signup endpoint"))
	}
}

func (s *Server) handleLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login endpoint"))
	}
}
