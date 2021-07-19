package database

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func BeforeGroupTest() (*DBHandler, int, error) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		return nil, -1, err
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		return nil, -1, err
	}
	user, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		return nil, -1, err
	}
	return handler, user.ID, nil
}

func AfterGroupTest(handler *DBHandler, id int) {
	_ = handler.DeleteUser(id)
	_ = handler.Close()
}

func TestAddAndDeleteGroup(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	group, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	if group.Name != "TestGroup" {
		t.Errorf("Failed to affect the name of the group")
	}
	if group.UserID != userID {
		t.Errorf("Failed to affect the user id of the group")
	}
	err = handler.DeleteGroup(group.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestAddGroupWrongUserID(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	_, err = handler.AddGroup(-1, "TestGroup")
	if err == nil {
		t.Errorf("Failed to add wrong user id : use a wrong user id didn't return an error")
	}
	AfterGroupTest(handler, userID)
}

func TestModifyGroupRight(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.ModifyGroup(tmpGroup.ID, "Toto")
	if err != nil {
		t.Errorf("Failed to modify the group : " + err.Error())
	}
	if group.ID != tmpGroup.ID {
		t.Errorf("Failed to modify group name, the group modified is different of the group provided")
	}
	if group.Name != "Toto" {
		t.Errorf("Failed to modify group, the group name has not been modified correctly")
	}
	if group.UserID != tmpGroup.UserID {
		t.Errorf("Failed to modify group, the group user ID have changed")
	}
	err = handler.DeleteGroup(group.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestModifyGroupFalseID(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.ModifyGroup(-1, "Toto")
	if err != nil {
		t.Errorf("Failed to modify group : " + err.Error())
	}
	if group != nil {
		t.Errorf("Failed to modify group : A wrong ID returns a group")
	}
	err = handler.DeleteGroup(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestGetGroupByIDRight(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.GetGroupByID(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to get group by its ID : " + err.Error())
	}
	if group.ID != tmpGroup.ID {
		t.Errorf("Failed to get group by its ID : The id of the group have changed")
	}
	if group.Name != tmpGroup.Name {
		t.Errorf("Failed to get group by its ID : The name of the group have changed")
	}
	if group.UserID != tmpGroup.UserID {
		t.Errorf("Failed to get group by its ID : The user owner of the group have changed")
	}
	err = handler.DeleteGroup(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestGetGroupByIDWrong(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	group, err := handler.GetGroupByID(-1)
	if err != nil {
		t.Errorf("Failed to get group by its ID : " + err.Error())
	}
	if group != nil {
		t.Errorf("Failed to get group by its ID : A wrong id returns a result")
	}
	AfterGroupTest(handler, userID)
}

func TestGetGroupByNameUserRight(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.GetGroupByUserAndName(tmpGroup.UserID, tmpGroup.Name)
	if err != nil {
		t.Errorf("Failed to get group by its user and its name : " + err.Error())
	}
	if group.ID != tmpGroup.ID {
		t.Errorf("Failed to get group by its user and its name : The id of the group have changed")
	}
	if group.Name != tmpGroup.Name {
		t.Errorf("Failed to get group by its user and its name : The name of the group have changed")
	}
	if group.UserID != tmpGroup.UserID {
		t.Errorf("Failed to get group by its user and its name : The user owner of the group have changed")
	}
	err = handler.DeleteGroup(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestGetGroupByUserAndNameWrongUser(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.GetGroupByUserAndName(-1, tmpGroup.Name)
	if err != nil {
		t.Errorf("Failed to get group by its user and name : " + err.Error())
	}
	if group != nil {
		t.Errorf("Failed to get group by its user and name : A wrong user returns a result")
	}
	err = handler.DeleteGroup(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}

func TestGetGroupByUserAndNameWrongName(t *testing.T) {
	handler, userID, err := BeforeGroupTest()
	if err != nil {
		t.Errorf("Failed to initialize test group, see others tests to find more explanation : " + err.Error())
	}
	tmpGroup, err := handler.AddGroup(userID, "TestGroup")
	if err != nil {
		t.Errorf("Failed to add group : " + err.Error())
	}
	group, err := handler.GetGroupByUserAndName(tmpGroup.ID, "Toto")
	if err != nil {
		t.Errorf("Failed to get group by its user and name : " + err.Error())
	}
	if group != nil {
		t.Errorf("Failed to get group by its user and name : A wrong name returns a result")
	}
	err = handler.DeleteGroup(tmpGroup.ID)
	if err != nil {
		t.Errorf("Failed to delete the group : " + err.Error())
	}
	AfterGroupTest(handler, userID)
}
