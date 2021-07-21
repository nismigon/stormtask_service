package web

import (
	"bytes"
	"encoding/json"
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
	found := false
	cookies := response.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Token" && cookie.Value != "" {
			found = true
		}
	}
	if !found {
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
