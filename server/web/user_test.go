package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"teissem/stormtask/server/configuration"
	"teissem/stormtask/server/database"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	Server       *Server
	HTTPServer   *httptest.Server
	User         *database.UserInformation
	UserPassword string
}

func (suite *UserTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		suite.T().Errorf("Failed to init the web server : " + err.Error())
	}
	user, err := server.Database.AddUser("web_user_test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	httpServer := httptest.NewServer(server.Router)
	suite.Server = server
	suite.HTTPServer = httpServer
	suite.User = user
	suite.UserPassword = "Test"
}

func (suite *UserTestSuite) TearDownTest() {
	_ = suite.Server.Database.DeleteUser(suite.User.ID)
	suite.HTTPServer.Close()
	_ = suite.Server.Close()
}

func (suite *UserTestSuite) TestAuthenticateRight() {
	var cred Credentials
	cred.Email = suite.User.Email
	cred.Password = suite.UserPassword
	authJSON, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to authenticate the user")
	}
	cookie := GetCookieByNameForResponse(response, suite.Server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to set the token in the cookie")
	}
}

func (suite *UserTestSuite) TestAuthenticateWrongEmail() {
	var cred Credentials
	cred.Email = "toto@toto.com"
	cred.Password = suite.User.Password
	authJSON, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf(
			"Failed to authenticate the user : A wrong email need to be reply with a unauthorized HTTP code")
	}
}

func (suite *UserTestSuite) TestAuthenticateWrongPassword() {
	var cred Credentials
	cred.Email = suite.User.Email
	cred.Password = "Toto"
	authJSON, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		suite.T().Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf(
			"Failed to authenticate the user : A wrong email need to be reply with a unauthorized HTTP code")
	}
}

func (suite *UserTestSuite) TestAddUserRight() {
	user := UserBody{
		Email:    "testAddUser@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err := json.Marshal(user)
	if err != nil {
		suite.T().Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/user", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the data on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to add the user with error code : " + fmt.Sprint(response.StatusCode))
	}
	userInfo, err := suite.Server.Database.GetUserByEmail("testAddUser@test.com")
	if err != nil {
		suite.T().Errorf("Failed to get the added user : " + err.Error())
	}
	_ = suite.Server.Database.DeleteUser(userInfo.ID)
}

func (suite *UserTestSuite) TestAddUserWrongAlreadyExistEmail() {
	userStruct := UserBody{
		Email:    "web_user_test@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err := json.Marshal(userStruct)
	if err != nil {
		suite.T().Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/user", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the data on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf(
			"Failed to add the user with error code : A dupplicated email doesn't return an unauthorized code")
	}
}

func (suite *UserTestSuite) TestDeleteUserRight() {
	cred := Credentials{
		Email:    "web_user_test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, suite.Server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", suite.HTTPServer.URL+"/user", nil)
	if err != nil {
		suite.T().Errorf("Failed to create a delete request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to delete the user with error code : Should retun a status 200")
	}
	deletedUser, err := suite.Server.Database.GetUserByID(suite.User.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the user by it's id : " + err.Error())
	}
	if deletedUser != nil {
		suite.T().Errorf("The user have not been deleted")
	}
}

func (suite *UserTestSuite) TestDeleteUserWrongToken() {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", suite.HTTPServer.URL+"/user", nil)
	if err != nil {
		suite.T().Errorf("Failed to create a delete request : " + err.Error())
	}
	req.Header.Set("Cookie", suite.Server.Configuration.TokenCookieName+"=Toto")
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf("Failed to delete the user with error code : Should retun a status 401")
	}
}

func (suite *UserTestSuite) TestModifyUserRight() {
	cred := Credentials{
		Email:    "web_user_test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, suite.Server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "test2@test.com",
		Password: "Test2",
		Name:     "Test2",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		suite.T().Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to delete the user with error code : Should retun a status 200")
	}
}

func (suite *UserTestSuite) TestModifyUserWrongAlreadyTakenEmail() {
	userForError, err := suite.Server.Database.AddUser("test2@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	cred := Credentials{
		Email:    "web_user_test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, suite.Server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "test2@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		suite.T().Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusConflict {
		suite.T().Errorf("Failed to delete the user with error code : Should retun a status 409")
	}
	err = suite.Server.Database.DeleteUser(userForError.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user from the database : " + err.Error())
	}
}

func (suite *UserTestSuite) TestModifyUserWrongInvalidToken() {
	cred := Credentials{
		Email:    "web_user_test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		suite.T().Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(suite.HTTPServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		suite.T().Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, suite.Server.Configuration.TokenCookieName)
	if cookie == nil {
		suite.T().Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "web_user_test@test.com",
		Password: "Test2",
		Name:     "Test2",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		suite.T().Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", suite.HTTPServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		suite.T().Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"=MyValue")
	}
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		suite.T().Errorf("Failed to delete the user with error code : Should return a status 401")
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
