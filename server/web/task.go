package web

import (
	"encoding/json"
	"net/http"
)

type TaskWithoutIDBody struct {
	Name        string
	Description string
	IsFinished  bool
	IsArchived  bool
	IDGroup     int
}

type TaskCompleteBody struct {
	ID          int
	Name        string
	Description string
	IsFinished  bool
	IsArchived  bool
	IDGroup     int
}

type TaskIDBody struct {
	ID int
}

type TaskGroupIDBody struct {
	GroupID int
}

// AddTask add a task for the authenticated user
// In the nominal case, this function returns a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not found, this returns a 404 HTTP code (Not found)
// If the group of the task is not owned by the user, this function returns a 401 HTTP code (Unauthorized)
// If the task is not added to the database, this function returns a 401 HTTP code (Unauthorized)
func (s *Server) AddTask(w http.ResponseWriter, r *http.Request) {
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
	var taskBody TaskWithoutIDBody
	err = json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	group, err := s.Database.GetGroupByID(taskBody.IDGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = s.Database.AddTask(taskBody.Name, taskBody.Description, taskBody.IsArchived,
		taskBody.IsArchived, taskBody.IDGroup)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ModifyTask modifies a task
// In the nominal case, this functions a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not found, this returns a 404 HTTP code (Not found)
// If the group of the task is not owned by the user, this function returns a 401 HTTP code (Unauthorized)
// If the future group of the task is not owned by the user, this function returns a 401 HTTP code (Unauthorized)
// If the task is not modified to the database, this function returns a 401 HTTP code (Unauthorized)
func (s *Server) ModifyTask(w http.ResponseWriter, r *http.Request) {
	// Verify and get the token in the cookie
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
	var taskBody TaskCompleteBody
	err = json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Verify if the group is owned by the selected user
	group, err := s.Database.GetGroupByID(taskBody.IDGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Verify if the task exists
	task, err := s.Database.GetTaskByID(taskBody.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = s.Database.ModifyTask(taskBody.ID, taskBody.Name, taskBody.Description, taskBody.IsArchived,
		taskBody.IsArchived, taskBody.IDGroup)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteTask deletes a task owned by the user
// In the nominal case, this functions a 200 HTTP code (OK)
// If the cookie is not found, or if the token is invalid, this function returns a 401 HTTP code (Unauthorized)
// If the body can be parsed, this function returns a 400 HTTP code (Bad request)
// If the group is not found, this returns a 404 HTTP code (Not found)
// If the group of the task is not owned by the user, this function returns a 401 HTTP code (Unauthorized)
// If the task is not deleted into the database, this function returns a 401 HTTP code (Unauthorized)
func (s *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Verify and get the token in the cookie
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
	// Get the body of the request
	var taskBody TaskIDBody
	err = json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := s.Database.GetTaskByID(taskBody.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Verification if the group and if it's owned by the authenticated user
	group, err := s.Database.GetGroupByID(task.IDGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Delete the task
	err = s.Database.DeleteTask(taskBody.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Verify and get the token in the cookie
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
	var taskBody TaskGroupIDBody
	err = json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	group, err := s.Database.GetGroupByID(taskBody.GroupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if group == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if group.UserID != claims.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tasks, err := s.Database.GetTasksByGroup(taskBody.GroupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
