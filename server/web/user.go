package web

import (
	"encoding/json"
	"net/http"
)

// Authenticate is an HTTP handler method to get a client authentication
func (s *Server) Authenticate(rw http.ResponseWriter, r *http.Request) {
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := s.generateToken(cred)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if token == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.SetCookie(rw, &http.Cookie{
		Name:  "Token",
		Value: token,
	})
	rw.WriteHeader(http.StatusOK)
}
