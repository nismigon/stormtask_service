package database

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func TestAddUser(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
		t.FailNow()
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		t.Errorf("Failed to open the database : " + err.Error())
		t.FailNow()
	}
	defer handler.Close()
	user, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	if user.Email != "test@test.com" {
		t.Errorf("Failed to set the email address")
	}
	if user.Name != "Test" {
		t.Errorf("Failed to set the name of the user")
	}
	if user.IsAdmin {
		t.Errorf("Failed to set the admin status of the user")
	}
}
