package server

func (s *Server) Routes() {
	s.router.Get("/health", s.handleHealthGet())
	s.router.Get("/", s.handleHelloGet())
	s.router.Post("/signup", s.handleSignupPost())
	s.router.Post("/login", s.handleLoginPost())
	s.router.Get("/me", s.authMiddleware(s.handleMeGet()))
}
