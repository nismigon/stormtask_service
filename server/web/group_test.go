package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"naleakan/stormtask/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BeforeGroupTest() (*Server, *httptest.Server, int, int, *http.Cookie, error) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	server, err := InitServer(*conf)
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	group, err := server.Database.AddGroup(user.ID, "TestGroup")
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	httpServer := httptest.NewServer(server.Router)
	cred := Credentials{
		Email:    "test@test.com",
		Password: "Test",
	}
	content, err := json.Marshal(cred)
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	response, err := http.Post(httpServer.URL+"/authenticate", "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, nil, -1, -1, nil, err
	}
	if response.StatusCode != http.StatusOK {
		err := errors.New("Failed to authenticate the user with error code : Should return a status 200")
		return nil, nil, -1, -1, nil, err
	}
	cookie := GetCookieByNameForResponse(response, server.Configuration.TokenCookieName)
	if cookie == nil {
		err := errors.New("Failed to find the token cookie")
		return nil, nil, -1, -1, nil, err
	}
	return server, httpServer, user.ID, group.ID, cookie, nil
}

func AfterGroupTest(server *Server, httpServer *httptest.Server, groupID int) {
	group, _ := server.Database.GetGroupByID(groupID)
	_ = server.Database.DeleteGroup(group.ID)
	_ = server.Database.DeleteUser(group.UserID)
	httpServer.Close()
	_ = server.Close()
}

func TestAddGroupRight(t *testing.T) {
	server, httpServer, userID, groupID, cookie, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize group test, please other test to know what happens : " + err.Error())
	}
	groupBody := AddGroupBody{
		Name: "MyNameForANewGroup",
	}
	groupJSON, err := json.Marshal(groupBody)
	if err != nil {
		t.Errorf("Failed to convert into JSON the AddGroupBody object : " + err.Error())
	}
	body := bytes.NewBuffer(groupJSON)
	req, err := http.NewRequest("POST", httpServer.URL+"/group", body)
	if err != nil {
		t.Errorf("Failed to create a post request for group : " + err.Error())
	}
	if cookie != nil {
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to get the response for the get route : " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Failed to add the group : " + err.Error())
	}
	group, err := server.Database.GetGroupByUserAndName(userID, "MyNameForANewGroup")
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	if group == nil {
		t.Errorf("Failed to find the added group")
		AfterGroupTest(server, httpServer, groupID)
		return
	}
	err = server.Database.DeleteGroup(group.ID)
	if err != nil {
		t.Errorf("Failedr to delete the group : " + err.Error())
	}
	AfterGroupTest(server, httpServer, groupID)
}
