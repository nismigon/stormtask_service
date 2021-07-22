package web

func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/authenticate", s.Authenticate).Methods("POST")
	s.Router.HandleFunc("/user", s.AddUser).Methods("POST")
	s.Router.HandleFunc("/user", s.DeleteUser).Methods("DELETE")

}
