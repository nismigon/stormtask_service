package database

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	Handler *DBHandler
}

func (suite *UserTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
		suite.T().FailNow()
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		suite.T().Errorf("Failed to open the database : " + err.Error())
		suite.T().FailNow()
	}
	suite.Handler = handler
}

func (suite *UserTestSuite) TearDownTest() {
	_ = suite.Handler.Close()
}

func (suite *UserTestSuite) TestAddAndDeleteUser() {
	user, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	if user.Email != "test@test.com" {
		suite.T().Errorf("Failed to set the email address")
	}
	if user.Name != "Test" {
		suite.T().Errorf("Failed to set the name of the user")
	}
	if user.Password != "Test" {
		suite.T().Errorf("Failed tp set the password of the user")
	}
	if user.IsAdmin {
		suite.T().Errorf("Failed to set the admin status of the user")
	}
	err = suite.Handler.DeleteUser(user.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *UserTestSuite) TestAddTwiceSameEmail() {
	userFirst, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	userSecond, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err == nil {
		suite.T().Errorf("Failed : no errors have been created when user have been added twice in the database")
		err = suite.Handler.DeleteUser(userFirst.ID)
		if err != nil {
			suite.T().Errorf("Fail to delete the first user")
		}
		err = suite.Handler.DeleteUser(userSecond.ID)
		if err != nil {
			suite.T().Errorf("Fail to delete the second user")
		}
		suite.T().FailNow()
	}
	err = suite.Handler.DeleteUser(userFirst.ID)
	if err != nil {
		suite.T().Errorf("Fail to delete the user")
	}
}

func (suite *UserTestSuite) TestGetUserByIDValid() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	user, err := suite.Handler.GetUserByID(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the user by its ID : " + err.Error())
		suite.T().FailNow()
	}
	if user.ID != tmpUser.ID {
		suite.T().Errorf("Failed to get the user with the right ID")
	}
	if user.Email != tmpUser.Email {
		suite.T().Errorf("Failed to get the user with the right email")
	}
	if user.Name != tmpUser.Name {
		suite.T().Errorf("Failed to get the user with the right name")
	}
	if user.Password != tmpUser.Password {
		suite.T().Errorf("Failed to get the user with the right password")
	}
	if user.IsAdmin != tmpUser.IsAdmin {
		suite.T().Errorf("Failed to get the user with the right admin status")
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user")
	}
}

func (suite *UserTestSuite) TestAuthenticateRight() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	user, err := suite.Handler.Authenticate(tmpUser.Email, tmpUser.Password)
	if err != nil {
		suite.T().Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user == nil {
			suite.T().Errorf("Failed to authenticate the user in the database, no error so maybe a code problem")
		}
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *UserTestSuite) TestAuthenticateWrongEmail() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	user, err := suite.Handler.Authenticate("toto@test.com", tmpUser.Password)
	if err != nil {
		suite.T().Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user != nil {
			suite.T().Errorf("Failed to authenticate the user in the database : no error, whereas provided wrong email")
		}
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *UserTestSuite) TestAuthenticateWrongPassword() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	user, err := suite.Handler.Authenticate(tmpUser.Email, "Toto")
	if err != nil {
		suite.T().Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user != nil {
			suite.T().Errorf("Failed to authenticate the user in the database : no error, whereas provided wrong password")
		}
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *UserTestSuite) TestModifyUserRight() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	user, err := suite.Handler.ModifyUser(tmpUser.ID, "toto@toto.com", "Toto", "Toto")
	if err != nil {
		suite.T().Errorf("Failed to modify the user in the database : " + err.Error())
	} else {
		if user.Email != "toto@toto.com" {
			suite.T().Errorf("Failed to modify the email of the user in the database")
		}
		if user.Name != "Toto" {
			suite.T().Errorf("Failed to modify the name of the user in the database")
		}
		if user.Password != "Toto" {
			suite.T().Errorf("Failed to modify the password of the user in the database")
		}
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func (suite *UserTestSuite) TestDeleteUserWithGroups() {
	tmpUser, err := suite.Handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
		suite.T().FailNow()
	}
	_, err = suite.Handler.AddGroup(tmpUser.ID, "GroupTest")
	if err != nil {
		suite.T().Errorf("Failed to add the group into the database : " + err.Error())
	}
	err = suite.Handler.DeleteUser(tmpUser.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
