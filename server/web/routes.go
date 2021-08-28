package web

import "net/http"

func (s *Server) InitRoutes() {
	handlerAuthenticate := http.HandlerFunc(s.Authenticate)
	handlerAddUser := http.HandlerFunc(s.AddUser)
	handlerDeleteUser := http.HandlerFunc(s.DeleteUser)
	handlerModifyUser := http.HandlerFunc(s.ModifyUser)
	handlerAddGroup := http.HandlerFunc(s.AddGroup)
	handlerGetGroups := http.HandlerFunc(s.GetGroupsByUserID)
	handlerModifyGroup := http.HandlerFunc(s.ModifyGroup)
	s.Router.Handle("/authenticate", handlerAuthenticate).Methods("POST")
	s.Router.Handle("/user", handlerAddUser).Methods("POST")
	s.Router.Handle("/user", s.VerifyToken(handlerDeleteUser)).Methods("DELETE")
	s.Router.Handle("/user", s.VerifyToken(handlerModifyUser)).Methods("PUT")
	s.Router.Handle("/group", s.VerifyToken(handlerAddGroup)).Methods("POST")
	s.Router.Handle("/group", s.VerifyToken(handlerGetGroups)).Methods("GET")
	s.Router.Handle("/group", s.VerifyToken(handlerModifyGroup)).Methods("PUT")
}
