package web

import "net/http"

func (s *Server) InitRoutes() {
	handlerAuthenticate := http.HandlerFunc(s.Authenticate)
	handlerAddUser := http.HandlerFunc(s.AddUser)
	handlerDeleteUser := http.HandlerFunc(s.DeleteUser)
	s.Router.Handle("/authenticate", handlerAuthenticate).Methods("POST")
	s.Router.Handle("/user", handlerAddUser).Methods("POST")
	s.Router.Handle("/user", s.VerifyToken(handlerDeleteUser)).Methods("DELETE")
}
