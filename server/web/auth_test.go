package web

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func TestGenerateTokenRight(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		t.Errorf("Failed to init the server : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    user.Email,
		Password: user.Password,
	}
	token, err := server.generateToken(cred)
	if err != nil {
		t.Errorf("Failed to generate the token : " + err.Error())
	}
	if token == "" {
		t.Errorf("Failed to generate the token")
	}
	claims, err := server.validateAndExtractToken(token)
	if err != nil {
		t.Errorf("Failed to validate and extract the token : " + err.Error())
	}
	if claims == nil {
		t.Errorf("Failed to validate and extract the token : The token is invalid")
	} else {
		if claims.ID != user.ID {
			t.Errorf("Failed to assign the ID in the token")
		}
		if claims.Email != user.Email {
			t.Errorf("Failed to assign the email in the token")
		}
		if claims.Name != user.Name {
			t.Errorf("Failed to assign the name in the token")
		}
		if claims.IsAdmin != user.IsAdmin {
			t.Errorf("Failed to assign the is admin in the token")
		}
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestGenerateTokenWrongEmail(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		t.Errorf("Failed to init the server : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    "toto@toto.com",
		Password: user.Password,
	}
	token, err := server.generateToken(cred)
	if err != nil {
		t.Errorf("Failed to generate the token : " + err.Error())
	}
	if token != "" {
		t.Errorf("Failed to generate the token : A wrong email give a correct token")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestGenerateTokenWrongPassword(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		t.Errorf("Failed to init the server : " + err.Error())
	}
	user, err := server.Database.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    user.Email,
		Password: "toto",
	}
	token, err := server.generateToken(cred)
	if err != nil {
		t.Errorf("Failed to generate the token : " + err.Error())
	}
	if token != "" {
		t.Errorf("Failed to generate the token : A wrong password give a correct token")
	}
	err = server.Database.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}
