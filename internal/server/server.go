package server

import (
	"database/sql"
	"log"
	"net/http"

	"shareapp/internal/db"
	"shareapp/utils"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router   *chi.Mux
	port     string
	db       *sql.DB
	queries  *db.Queries
	jwtMaker *utils.JWTMaker
}

func New(port, dsn string) (*Server, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to db", "err", err)
	}

	s := &Server{
		router:   chi.NewRouter(),
		port:     port,
		db:       dbConn,
		queries:  db.New(dbConn),
		jwtMaker: utils.NewJWTMaker("secret-key"),
	}
	s.Routes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
