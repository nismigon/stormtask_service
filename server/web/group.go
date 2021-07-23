package web

import (
	"encoding/json"
	"net/http"
)

type AddGroupBody struct {
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
	var groupBody AddGroupBody
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
