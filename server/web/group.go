package web

import (
	"encoding/json"
	"net/http"
)

type GroupNameBody struct {
	Name string
}

type GroupIDNameBody struct {
	ID   int
	Name string
}

// AddGroup add a group for the authenticated user
// In the nominal case, this function returns a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not added to the database, this function returns a 401 HTTP code (Unauthorized)
func (server *Server) AddGroup(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookieByNameForRequest(r, server.Configuration.TokenCookieName)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, err := server.ValidateAndExtractToken(cookie.Value)
	if err != nil || claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var groupBody GroupNameBody
	err = json.NewDecoder(r.Body).Decode(&groupBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = server.Database.AddGroup(claims.ID, groupBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetGroupsByUserID get all the groups of a user
// In the nominal case, this returns a JSON object containing all the groups of the user
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If an error occurred when we try to get user groups or when encoding the JSON, this functions returns a 500
// HTTP code (Internal Server Error)
func (server *Server) GetGroupsByUserID(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookieByNameForRequest(r, server.Configuration.TokenCookieName)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, err := server.ValidateAndExtractToken(cookie.Value)
	if err != nil || claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	groups, err := server.Database.GetGroupsByUserID(claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
