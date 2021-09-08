package web

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	Server *Server
}

func (suite *AuthTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		suite.T().Errorf("Failed to init the server : " + err.Error())
	}
	suite.Server = server
}

func (suite *AuthTestSuite) TearDownTest() {
	_ = suite.Server.Close()
}

func (suite *AuthTestSuite) TestGenerateTokenRight() {
	user, err := suite.Server.Database.AddUser("web_auth_test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    user.Email,
		Password: "Test",
	}
	token, err := suite.Server.generateToken(cred)
	if err != nil {
		suite.T().Errorf("Failed to generate the token : " + err.Error())
	}
	if token == "" {
		suite.T().Errorf("Failed to generate the token")
	}
	claims, err := suite.Server.ValidateAndExtractToken(token)
	if err != nil {
		suite.T().Errorf("Failed to validate and extract the token : " + err.Error())
	}
	if claims == nil {
		suite.T().Errorf("Failed to validate and extract the token : The token is invalid")
	} else {
		if claims.ID != user.ID {
			suite.T().Errorf("Failed to assign the ID in the token")
		}
		if claims.Email != user.Email {
			suite.T().Errorf("Failed to assign the email in the token")
		}
		if claims.Name != user.Name {
			suite.T().Errorf("Failed to assign the name in the token")
		}
		if claims.IsAdmin != user.IsAdmin {
			suite.T().Errorf("Failed to assign the is admin in the token")
		}
	}
	err = suite.Server.Database.DeleteUser(user.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *AuthTestSuite) TestGenerateTokenWrongEmail() {
	user, err := suite.Server.Database.AddUser("web_auth_test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    "toto@toto.com",
		Password: user.Password,
	}
	token, err := suite.Server.generateToken(cred)
	if err != nil {
		suite.T().Errorf("Failed to generate the token : " + err.Error())
	}
	if token != "" {
		suite.T().Errorf("Failed to generate the token : A wrong email give a correct token")
	}
	err = suite.Server.Database.DeleteUser(user.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *AuthTestSuite) TestGenerateTokenWrongPassword() {
	user, err := suite.Server.Database.AddUser("web_auth_test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user : " + err.Error())
	}
	var cred = Credentials{
		Email:    user.Email,
		Password: "toto",
	}
	_, err = suite.Server.generateToken(cred)
	if err == nil {
		suite.T().Errorf("Failed to generate the token : No error found for a false credentials")
	}
	err = suite.Server.Database.DeleteUser(user.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
