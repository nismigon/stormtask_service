package web

import (
	"encoding/json"
	"net/http"
)

type GroupIDBody struct {
	ID int
}

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
func (s *Server) AddGroup(w http.ResponseWriter, r *http.Request) {
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
	var groupBody GroupNameBody
	err = json.NewDecoder(r.Body).Decode(&groupBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = s.Database.AddGroup(claims.ID, groupBody.Name)
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
func (s *Server) GetGroupsByUserID(w http.ResponseWriter, r *http.Request) {
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
	groups, err := s.Database.GetGroupsByUserID(claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ModifyGroup modifies a group of the user which own the token
// In the nominal case, this function returns a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not found, this function returns a 404 HTTP code (Not found)
// If the group is not owned by the user on the token, this function returns a 401 HTTP code (Unauthorized)
// If an error occurred during the change of the group name, this function returns a
// 500 HTTP code (Internal Server Error)
func (s *Server) ModifyGroup(w http.ResponseWriter, r *http.Request) {
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
	var groupBody GroupIDNameBody
	err = json.NewDecoder(r.Body).Decode(&groupBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Verification if the group is owned by the user
	group, err := s.Database.GetGroupByID(groupBody.ID)
	if err != nil || group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Modification of the group
	_, err = s.Database.ModifyGroup(groupBody.ID, groupBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteGroup deletes a group if the the user which own the token
// In the nominal case, this function returns a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not found, this function returns a 404 HTTP code (Not found)
// If the group is not owned by the user on the token, this function returns a 401 HTTP code (Unauthorized)
// If an error occurred during the removal, this function returns a
// 500 HTTP code (Internal Server Error)
func (s *Server) DeleteGroup(w http.ResponseWriter, r *http.Request) {
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
	var groupBody GroupIDBody
	err = json.NewDecoder(r.Body).Decode(&groupBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Verification if the group is owned by the user
	group, err := s.Database.GetGroupByID(groupBody.ID)
	if err != nil || group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Delete the group
	err = s.Database.DeleteGroup(group.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
