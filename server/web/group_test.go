package web

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"teissem/stormtask/server/configuration"
	"teissem/stormtask/server/database"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GroupTestSuite struct {
	suite.Suite
	Server     *Server
	HTTPServer *httptest.Server
	UserID     int
	GroupID    int
	Cookie     *http.Cookie
}

func (suite *GroupTestSuite) SetupTest() {
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
	suite.UserID = user.ID
	suite.GroupID = group.ID
	suite.Cookie = cookie
}

func (suite *GroupTestSuite) TearDownTest() {
	_ = suite.Server.Database.DeleteUser(suite.UserID)
	suite.HTTPServer.Close()
	_ = suite.Server.Close()
}

func (suite *GroupTestSuite) TestAddGroupRight() {
	groupBody := GroupNameBody{
		Name: "MyNameForANewGroup",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/group", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for group : " + err.Error())
	}
	if suite.Cookie != nil {
		req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to add the group : " + err.Error())
	}
	group, err := suite.Server.Database.GetGroupByUserAndName(suite.UserID, "MyNameForANewGroup")
	if err != nil {
		suite.T().Errorf("Failed to delete the group : " + err.Error())
	}
	if group == nil {
		suite.T().Errorf("Failed to find the added group")
	}
}

func (suite *GroupTestSuite) TestAddGroupWrongName() {
	groupBody := GroupNameBody{
		Name: "TestGroup",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("POST", suite.HTTPServer.URL+"/group", body)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for group : " + err.Error())
	}
	if suite.Cookie != nil {
		req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the post route : " + err.Error())
	}
	if response.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf(
			"Failed to add the group : Add a group with the same name should return an Unauthorized HTTP code")
	}
}

func (suite *GroupTestSuite) TestGetGroupsByUserIDRight() {
	req, err := http.NewRequest("GET", suite.HTTPServer.URL+"/group", nil)
	if err != nil {
		suite.T().Errorf("Failed to create a post request for group : " + err.Error())
	}
	if suite.Cookie != nil {
		req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the get route : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Error("Failed to get the differents groups : return an error code " +
			strconv.FormatInt(int64(response.StatusCode), 10))
	}
	var groups []database.GroupInformation
	err = json.NewDecoder(response.Body).Decode(&groups)
	if err != nil {
		suite.T().Errorf("Failed to decode the body of the response with : " + err.Error())
	}
	if len(groups) != 1 {
		suite.T().Errorf("Failed to get the correct result : the array of group should contain 1 element")
	}
}

func (suite *GroupTestSuite) TestModifyGroupRight() {
	groupBody := GroupIDNameBody{
		ID:   suite.GroupID,
		Name: "New Group Name",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/group", body)
	if err != nil {
		suite.T().Errorf("Failed to create the request : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the modify route : " + err.Error())
	}
	assert.Equal(suite.T(), 200, response.StatusCode)
	group, err := suite.Server.Database.GetGroupByID(suite.GroupID)
	if err != nil {
		suite.T().Errorf("Failed to get the group : " + err.Error())
	}
	assert.NotNil(suite.T(), group)
	assert.Equal(suite.T(), "New Group Name", group.Name)
}

func (suite *GroupTestSuite) TestModifyGroupWrongNotFound() {
	groupBody := GroupIDNameBody{
		ID:   -1,
		Name: "New Group Name",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/group", body)
	if err != nil {
		suite.T().Errorf("Failed to create the request : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the modify route : " + err.Error())
	}
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *GroupTestSuite) TestModifyGroupWrongUserID() {
	user, err := suite.Server.Database.AddUser("test2@test2.com", "Test", "Pass", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user to the database : " + err.Error())
	}
	defer func(Database *database.DBHandler, id int) {
		err := Database.DeleteUser(id)
		if err != nil {
			suite.T().Errorf("Failed to delete the user with error : " + err.Error())
		}
	}(suite.Server.Database, user.ID)
	group, err := suite.Server.Database.AddGroup(user.ID, "MyGroup")
	if err != nil {
		suite.T().Errorf("Failed to add group to the database : " + err.Error())
	}
	groupBody := GroupIDNameBody{
		ID:   group.ID,
		Name: "New Group Name",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		suite.T().Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/group", body)
	if err != nil {
		suite.T().Errorf("Failed to create the request : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Cookie.Name+"="+suite.Cookie.Value)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the modify route : " + err.Error())
	}
	assert.Equal(suite.T(), 401, response.StatusCode)
}

func TestGroup(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}
