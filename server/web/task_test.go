package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"teissem/stormtask/server/configuration"
	"teissem/stormtask/server/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TaskTestSuite struct {
	suite.Suite
	User       *database.UserInformation
	Group      *database.GroupInformation
	Task       *database.TaskInformation
	Server     *Server
	HTTPServer *httptest.Server
	Cookie     *http.Cookie
}

func (suite *TaskTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		suite.T().Errorf("Failed to init the web server : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	group, err := server.Database.AddGroup(user.ID, "TestGroup")
	if err != nil {
		suite.T().Errorf("Failed to add the group : " + err.Error())
	}
	task, err := server.Database.AddTask("MyTask", "Description", false, false, group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add the task : " + err.Error())
	}
	httpServer := httptest.NewServer(server.Router)
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the credentials : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to send the authentication request : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to send the authentication request : An error is returned")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to get the cookie")
	}
	suite.Server = server
	suite.HTTPServer = httpServer
	suite.User = user
	suite.Group = group
	suite.Task = task
	suite.Cookie = cookie
}

func (suite *TaskTestSuite) TearDownTest() {
	_ = suite.Server.Database.DeleteUser(suite.User.ID)
	suite.HTTPServer.Close()
	_ = suite.Server.Close()
}

func (suite *TaskTestSuite) TestAddTaskRight() {
	taskBody := TaskWithoutIDBody{
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     suite.Group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskWithoutIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 200, response.StatusCode)
	tasks, err := suite.Server.Database.GetTasksByGroup(suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the tasks : " + err.Error())
	}
	assert.Equal(suite.T(), 2, len(*tasks))
}

func (suite *TaskTestSuite) TestAddTaskWrongGroupForUser() {
	user, err := suite.Server.Database.AddUser("user@user.com", "User", "Password", false)
	assert.Nil(suite.T(), err)
	defer func(Database *database.DBHandler, id int) {
		err := Database.DeleteUser(id)
		if err != nil {
			suite.T().Errorf("Failed to delete the user with error : " + err.Error())
		}
	}(suite.Server.Database, user.ID)
	group, err := suite.Server.Database.AddGroup(user.ID, "MyGroup")
	assert.Nil(suite.T(), err)
	taskBody := TaskWithoutIDBody{
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskWithoutIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 401, response.StatusCode)
}

func (suite *TaskTestSuite) TestAddTaskWrongGroupId() {
	taskBody := TaskWithoutIDBody{
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     -1,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskWithoutIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *TaskTestSuite) TestModifyTaskRight() {
	taskBody := TaskCompleteBody{
		ID:          suite.Task.ID,
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     suite.Group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskCompleteBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 200, response.StatusCode)
	task, err := suite.Server.Database.GetTaskByID(suite.Task.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the task : " + err.Error())
	}
	assert.Equal(suite.T(), "My Task", task.Name)
	assert.Equal(suite.T(), "Description of my task", task.Description)
	assert.Equal(suite.T(), false, task.IsFinished)
	assert.Equal(suite.T(), false, task.IsArchived)
	assert.Equal(suite.T(), suite.Group.ID, task.IDGroup)
}

func (suite *TaskTestSuite) TestModifyTaskWrongGroupForUser() {
	user, err := suite.Server.Database.AddUser("user@user.com", "User", "Password", false)
	assert.Nil(suite.T(), err)
	defer func(Database *database.DBHandler, id int) {
		err := Database.DeleteUser(id)
		if err != nil {
			suite.T().Errorf("Failed to delete the user with error : " + err.Error())
		}
	}(suite.Server.Database, user.ID)
	group, err := suite.Server.Database.AddGroup(user.ID, "MyGroup")
	assert.Nil(suite.T(), err)
	taskBody := TaskCompleteBody{
		ID:          suite.Task.ID,
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskCompleteBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 401, response.StatusCode)
}

func (suite *TaskTestSuite) TestModifyTaskWrongGroupId() {
	taskBody := TaskCompleteBody{
		ID:          suite.Task.ID,
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     -1,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskWithoutIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *TaskTestSuite) TestModifyTaskWrongTask() {
	taskBody := TaskCompleteBody{
		ID:          -1,
		Name:        "My Task",
		Description: "Description of my task",
		IsFinished:  false,
		IsArchived:  false,
		IDGroup:     suite.Group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskCompleteBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *TaskTestSuite) TestDeleteTaskRight() {
	taskBody := TaskIDBody{
		ID: suite.Task.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("DELETE", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the delete route : " + err.Error())
	}
	assert.Equal(suite.T(), 200, response.StatusCode)
	task, err := suite.Server.Database.GetTaskByID(suite.Task.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the task : " + err.Error())
	}
	assert.Nil(suite.T(), task)
}

func (suite *TaskTestSuite) TestGetTasksRight() {
	taskBody := TaskGroupIDBody{
		GroupID: suite.Group.ID,
	}
	taskJSON, err := json.Marshal(taskBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the TaskGroupIDBody object : " + err.Error())
	}
	body := bytes.NewBuffer(taskJSON)
	req, err := http.NewRequest("GET", suite.HTTPServer.URL+"/task", body)
	if err != nil {
		suite.T().Errorf("Failed to create a get request for task : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	assert.Equal(suite.T(), 200, response.StatusCode)
	var tasks []database.TaskInformation
	err = json.NewDecoder(response.Body).Decode(&tasks)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Equal(suite.T(), suite.Task.ID, tasks[0].ID)
}

func TestTask(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}
