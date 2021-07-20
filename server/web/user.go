package web

import (
	"encoding/json"
	"net/http"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Authenticate is an HTTP handler method to get a client authentication
func (s *Server) Authenticate(rw http.ResponseWriter, r *http.Request) {
	var auth Authentication
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	user, err := s.Database.Authenticate(auth.Email, auth.Password)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	if user == nil {
		rw.WriteHeader(http.StatusUnauthorized)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
