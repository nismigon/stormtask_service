package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"naleakan/stormtask/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BeforeUserTest() (*Server, *httptest.Server, error) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		return nil, nil, err
	}
	server, err := InitServer(*conf)
	if err != nil {
		return nil, nil, err
	}
	httpServer := httptest.NewServer(server.Router)
	return server, httpServer, nil
}

func AfterUserTest(server *Server, httpServer *httptest.Server) {
	httpServer.Close()
	_ = server.Close()
}

func TestAuthenticateRight(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred Credentials
	cred.Email = user.Email
	cred.Password = user.Password
	authJSON, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		t.Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to authenticate the user")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		t.Errorf("Failed to set the token in the cookie")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestAuthenticateWrongEmail(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred Credentials
	cred.Email = "toto@toto.com"
	cred.Password = user.Password
	authJSON, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		t.Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("Failed to authenticate the user : A wrong email need to be reply with a unauthorized HTTP code")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestAuthenticateWrongPassword(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred Credentials
	cred.Email = user.Email
	cred.Password = "Toto"
	authJSON, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to convert the authencication object into JSON : " + err.Error())
	}
	body := bytes.NewBuffer(authJSON)
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", body)
	if err != nil {
		t.Errorf("Failed to get the response for the authenticate route : " + err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("Failed to authenticate the user : A wrong email need to be reply with a unauthorized HTTP code")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestAddUserRight(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user := UserBody{
		Email:    "test@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/user", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the data on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to add the user with error code : " + fmt.Sprint(response.StatusCode))
	}
	userInfo, err := server.Database.GetUserByEmail("test@test.com")
	if err != nil {
		t.Errorf("Failed to get the added user : " + err.Error())
	}
	err = server.Database.DeleteUser(userInfo.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestAddUserWrongAlreadyExistEmail(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	_, err = server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	userStruct := UserBody{
		Email:    "test@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err := json.Marshal(userStruct)
	if err != nil {
		t.Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/user", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the data on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("Failed to add the user with error code : A dupplicated email doesn't return an unauthorized code")
	}
	userInfo, err := server.Database.GetUserByEmail("test@test.com")
	if err != nil {
		t.Errorf("Failed to get the added user : " + err.Error())
	}
	err = server.Database.DeleteUser(userInfo.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestDeleteUserRight(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		t.Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", httpServer.URL+"/user", nil)
	if err != nil {
		t.Errorf("Failed to create a delete request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to delete the user with error code : Should retun a status 200")
	}
	deletedUser, err := server.Database.GetUserByID(user.ID)
	if err != nil {
		t.Errorf("Failed to get the user by it's id : " + err.Error())
	}
	if deletedUser != nil {
		t.Errorf("The user have not been deleted")
		err = server.Database.DeleteUser(user.ID)
		if err != nil {
			t.Errorf("Failed to delete the user from the database : " + err.Error())
		}
	}
	AfterUserTest(server, httpServer)
}

func TestDeleteUserWrongToken(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", httpServer.URL+"/user", nil)
	if err != nil {
		t.Errorf("Failed to create a delete request : " + err.Error())
	}
	req.Header.Set("Cookie", server.Configuration.TokenCookieName+"=Toto")
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Failed to delete the user with error code : Should retun a status 401")
	}
	AfterUserTest(server, httpServer)
}

func TestModifyUserRight(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		t.Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "test2@test.com",
		Password: "Test2",
		Name:     "Test2",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		t.Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", httpServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to delete the user with error code : Should retun a status 200")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user from the database : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestModifyUserWrongAlreadyTakenEmail(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	userForError, err := server.Database.AddUser("test2@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		t.Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "test2@test.com",
		Password: "Test",
		Name:     "Test",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		t.Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", httpServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Failed to delete the user with error code : Should retun a status 409")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user from the database : " + err.Error())
	}
	err = server.Database.DeleteUser(userForError.ID)
	if err != nil {
		t.Errorf("Failed to delete the user from the database : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}

func TestModifyUserWrongInvalidToken(t *testing.T) {
	server, httpServer, err := BeforeUserTest()
	if err != nil {
		t.Errorf("Failed to initialize user test, please other test to know what happens : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		t.Errorf("Failed to marshal credentials : " + err.Error())
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to POST the request on the web server : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to authenticate the user with error code : Should return a status 200")
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		t.Errorf("Failed to find the token cookie")
	}
	client := &http.Client{}
	userStruct := UserBody{
		Email:    "test@test.com",
		Password: "Test2",
		Name:     "Test2",
	}
	content, err = json.Marshal(userStruct)
	if err != nil {
		t.Errorf("Failed to marshal the add user struct needed to execute the request : " + err.Error())
	}
	req, err := http.NewRequest("PUT", httpServer.URL+"/user", bytes.NewReader(content))
	if err != nil {
		t.Errorf("Failed to create a put request : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"=MyValue")
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to execute the HTTP request : " + err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Failed to delete the user with error code : Should return a status 401")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user from the database : " + err.Error())
	}
	AfterUserTest(server, httpServer)
}
