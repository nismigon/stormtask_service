package web

func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/authenticate", s.Authenticate).Methods("POST")
}
