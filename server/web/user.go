package web

import (
	"encoding/json"
	"net/http"
)

type UserBody struct {
	Email    string
	Password string
	Name     string
}

// Authenticate is an HTTP handler method to get a client authentication
// In the nominal case, this returns a 200 HTTP code (OK) and generate a cookie named token with the JWT
// If the body doesn't correspond to the Credential struct, this returns a 400 HTTP code (Bad Request)
// If an error occurred during the creation of the token, this returns a 500 HTTP code (Internal Server Error)
// If the returned token is an empty string, this returns a 401 HTTP code (Unauthorized)
func (s *Server) Authenticate(rw http.ResponseWriter, r *http.Request) {
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := s.generateToken(cred)
	if err != nil || token == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.SetCookie(rw, &http.Cookie{
		Name:  s.Configuration.TokenCookieName,
		Value: token,
	})
	rw.WriteHeader(http.StatusOK)
}

// AddUser add a user in the database
// In the nominal case, this returns a 200 HTTP code (OK)
// If the request is malformed, this returns a 400 HTTP code (Bad Request)
// If an error occurred when adding the user in the database, this returns a 401 HTTP code (Unauthorized) because
// it's probably a duplicated email
func (s *Server) AddUser(rw http.ResponseWriter, r *http.Request) {
	var userBody UserBody
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = s.Database.AddUser(userBody.Email, userBody.Name, userBody.Password, false)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// DeleteUser delete the user which own the token
// In the nominal case, this returns a 200 HTTP code (OK)
// If the token is not found or if the token is invalid, this returns a 401 HTTP code (Unauthorized)
// If the user is not deleted from the database, this returns a 500 HTTP code (Internal Server Error)
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookieByNameForRequest(r, s.Configuration.TokenCookieName)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, err := s.ValidateAndExtractToken(cookie.Value)
	if err != nil || claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = s.Database.DeleteUser(claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ModifyUser modify the user into the database
// In the nominal case, this function returns a 200 HTTP code (OK)
// If the cookie is not found or if the token is not valid, this function returns a 401 HTTP code (Unauthorized)
// If the modification of the user (without email) failed, this function returns a 500 HTTP code (Internal Server Error)
// If the modification of the user (email) failed, this function returns a 409 HTTP code (Conflict)
func (s *Server) ModifyUser(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookieByNameForRequest(r, s.Configuration.TokenCookieName)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, err := s.ValidateAndExtractToken(cookie.Value)
	if err != nil || claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var userBody UserBody
	err = json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Request to modify user informations except email
	user, err := s.Database.ModifyUser(claims.ID, claims.Email, userBody.Name, userBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Request to modify the user email
	if userBody.Email != claims.Email {
		_, err = s.Database.ModifyUser(claims.ID, userBody.Email, user.Name, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
