package database

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GroupTestSuite struct {
	suite.Suite
	Handler *DBHandler
	UserID  int
	Group   *GroupInformation
}

func (suite *GroupTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		suite.T().Errorf("Failed to initialize the connection with the database : " + err.Error())
	}
	suite.Handler = handler
	user, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user into the database : " + err.Error())
	}
	suite.UserID = user.ID
	group, err := suite.Handler.AddGroup(suite.UserID, "TestGroup")
	if err != nil {
		suite.T().Errorf("Failed to add group : " + err.Error())
	}
	suite.Group = group
}

func (suite *GroupTestSuite) TearDownTest() {
	_ = suite.Handler.DeleteUser(suite.UserID)
	_ = suite.Handler.Close()
}

func (suite *GroupTestSuite) TestAddAndDeleteGroup() {
	group, err := suite.Handler.AddGroup(suite.UserID, "TestGroupUnique")
	if err != nil {
		suite.T().Errorf("Failed to add group : " + err.Error())
	}
	if group.Name != "TestGroupUnique" {
		suite.T().Errorf("Failed to affect the name of the group")
	}
	if group.UserID != suite.UserID {
		suite.T().Errorf("Failed to affect the user id of the group")
	}
	err = suite.Handler.DeleteGroup(group.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the group : " + err.Error())
	}
}

func (suite *GroupTestSuite) TestAddGroupWrongUserID() {
	_, err := suite.Handler.AddGroup(-1, "TestGroup")
	if err == nil {
		suite.T().Errorf("Failed to add wrong user id : use a wrong user id didn't return an error")
	}
}

func (suite *GroupTestSuite) TestModifyGroupRight() {
	group, err := suite.Handler.ModifyGroup(suite.Group.ID, "Toto")
	if err != nil {
		suite.T().Errorf("Failed to modify the group : " + err.Error())
	}
	if group.ID != suite.Group.ID {
		suite.T().Errorf("Failed to modify group name, the group modified is different of the group provided")
	}
	if group.Name != "Toto" {
		suite.T().Errorf("Failed to modify group, the group name has not been modified correctly")
	}
	if group.UserID != suite.Group.UserID {
		suite.T().Errorf("Failed to modify group, the group user ID have changed")
	}
}

func (suite *GroupTestSuite) TestModifyGroupFalseID() {
	group, err := suite.Handler.ModifyGroup(-1, "Toto")
	if err != nil {
		suite.T().Errorf("Failed to modify group : " + err.Error())
	}
	if group != nil {
		suite.T().Errorf("Failed to modify group : A wrong ID returns a group")
	}
}

func (suite *GroupTestSuite) TestGetGroupByIDRight() {
	group, err := suite.Handler.GetGroupByID(suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to get group by its ID : " + err.Error())
	}
	if group.ID != suite.Group.ID {
		suite.T().Errorf("Failed to get group by its ID : The id of the group have changed")
	}
	if group.Name != suite.Group.Name {
		suite.T().Errorf("Failed to get group by its ID : The name of the group have changed")
	}
	if group.UserID != suite.Group.UserID {
		suite.T().Errorf("Failed to get group by its ID : The user owner of the group have changed")
	}
}

func (suite *GroupTestSuite) TestGetGroupByIDWrong() {
	group, err := suite.Handler.GetGroupByID(-1)
	if err != nil {
		suite.T().Errorf("Failed to get group by its ID : " + err.Error())
	}
	if group != nil {
		suite.T().Errorf("Failed to get group by its ID : A wrong id returns a result")
	}
}

func (suite *GroupTestSuite) TestGetGroupByNameUserRight() {
	group, err := suite.Handler.GetGroupByUserAndName(suite.Group.UserID, suite.Group.Name)
	if err != nil {
		suite.T().Errorf("Failed to get group by its user and its name : " + err.Error())
	}
	if group.ID != suite.Group.ID {
		suite.T().Errorf("Failed to get group by its user and its name : The id of the group have changed")
	}
	if group.Name != suite.Group.Name {
		suite.T().Errorf("Failed to get group by its user and its name : The name of the group have changed")
	}
	if group.UserID != suite.Group.UserID {
		suite.T().Errorf("Failed to get group by its user and its name : The user owner of the group have changed")
	}
}

func (suite *GroupTestSuite) TestGetGroupByUserAndNameWrongUser() {
	group, err := suite.Handler.GetGroupByUserAndName(-1, suite.Group.Name)
	if err != nil {
		suite.T().Errorf("Failed to get group by its user and name : " + err.Error())
	}
	if group != nil {
		suite.T().Errorf("Failed to get group by its user and name : A wrong user returns a result")
	}
}

func (suite *GroupTestSuite) TestGetGroupByUserAndNameWrongName() {
	group, err := suite.Handler.GetGroupByUserAndName(suite.Group.UserID, "Toto")
	if err != nil {
		suite.T().Errorf("Failed to get group by its user and name : " + err.Error())
	}
	if group != nil {
		suite.T().Errorf("Failed to get group by its user and name : A wrong name returns a result")
	}
}

func (suite *GroupTestSuite) TestUniqueGroupNameByUser() {
	group, err := suite.Handler.AddGroup(suite.Group.UserID, "TestGroup")
	if err == nil {
		suite.T().Errorf("Failed to add group : no error whereas break unique constraint on name")
		err = suite.Handler.DeleteGroup(group.ID)
		if err != nil {
			suite.T().Errorf("Failed to delete group : " + err.Error())
		}
	}
}

func (suite *GroupTestSuite) TestGetGroupsByUserIDRight() {
	group1, err := suite.Handler.AddGroup(suite.UserID, "TestGroup1")
	if err != nil {
		suite.T().Errorf("Failed to add the group : " + err.Error())
	}
	group2, err := suite.Handler.AddGroup(suite.UserID, "TestGroup2")
	if err != nil {
		suite.T().Errorf("Failed to add the group : " + err.Error())
	}
	groups, err := suite.Handler.GetGroupsByUserID(suite.UserID)
	if err != nil {
		suite.T().Errorf("Failed to get groups by user id : " + err.Error())
	}
	if len(*groups) != 3 {
		suite.T().Errorf("Failed to get groups by user id : The table returned should contain 2 elements")
	}
	for _, group := range *groups {
		if group != *group1 && group != *group2 && group != *suite.Group {
			suite.T().Errorf("Failed to get groups by user id : Unrocognized group")
		}
	}
}

func (suite *GroupTestSuite) TestGetGroupsByUserIDWrongUserID() {
	groups, err := suite.Handler.GetGroupsByUserID(-1)
	if err != nil {
		suite.T().Errorf("Failed to get groups by user id : " + err.Error())
	}
	if len(*groups) != 0 {
		suite.T().Errorf("Failed to get groups by user id : The len of the slice should be equal to 0")
	}
}

func (suite *GroupTestSuite) TestDeleteGroupsByUserRight() {
	_, err := suite.Handler.AddGroup(suite.UserID, "Test1")
	if err != nil {
		suite.T().Errorf("Failed to add group, see others test to fine more explanation : " + err.Error())
	}
	_, err = suite.Handler.AddGroup(suite.UserID, "Test2")
	if err != nil {
		suite.T().Errorf("Failed to add group, see others test to fine more explanation : " + err.Error())
	}
	err = suite.Handler.DeleteGroupsByUser(suite.UserID)
	if err != nil {
		suite.T().Errorf("Failed to delete groups by user : " + err.Error())
	}
	groups, err := suite.Handler.GetGroupsByUserID(suite.UserID)
	if err != nil {
		suite.T().Errorf("Failed to get groups by user id : " + err.Error())
	}
	if len(*groups) != 0 {
		suite.T().Errorf("Failed to get groups by user id : The len of the slice should be equal to 0")
	}
}

func TestGroup(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}
