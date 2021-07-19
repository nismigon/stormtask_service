package database

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func TestAddAndDeleteUser(t *testing.T) {
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
	if user.Password != "Test" {
		t.Errorf("Failed tp set the password of the user")
	}
	if user.IsAdmin {
		t.Errorf("Failed to set the admin status of the user")
	}
	err = handler.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestAddTwiceSameEmail(t *testing.T) {
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
	userFirst, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	userSecond, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err == nil {
		t.Errorf("Failed : no errors have been created when user have been added twice in the database")
		err = handler.DeleteUser(userFirst.ID)
		if err != nil {
			t.Errorf("Fail to delete the first user")
		}
		err = handler.DeleteUser(userSecond.ID)
		if err != nil {
			t.Errorf("Fail to delete the second user")
		}
		t.FailNow()
	}
	err = handler.DeleteUser(userFirst.ID)
	if err != nil {
		t.Errorf("Fail to delete the user")
	}
}

func TestGetUserByIDValid(t *testing.T) {
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
	tmpUser, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	user, err := handler.GetUserByID(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to get the user by its ID : " + err.Error())
		t.FailNow()
	}
	if user.ID != tmpUser.ID {
		t.Errorf("Failed to get the user with the right ID")
	}
	if user.Email != tmpUser.Email {
		t.Errorf("Failed to get the user with the right email")
	}
	if user.Name != tmpUser.Name {
		t.Errorf("Failed to get the user with the right name")
	}
	if user.Password != tmpUser.Password {
		t.Errorf("Failed to get the user with the right password")
	}
	if user.IsAdmin != tmpUser.IsAdmin {
		t.Errorf("Failed to get the user with the right admin status")
	}
	err = handler.DeleteUser(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to delete the user")
	}
}

func TestAuthenticateRight(t *testing.T) {
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
	tmpUser, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	user, err := handler.Authenticate(tmpUser.Email, tmpUser.Password)
	if err != nil {
		t.Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user == nil {
			t.Errorf("Failed to authenticate the user in the database, no error so maybe a code problem")
		}
	}
	err = handler.DeleteUser(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestAuthenticateWrongEmail(t *testing.T) {
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
	tmpUser, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	user, err := handler.Authenticate("toto@test.com", tmpUser.Password)
	if err != nil {
		t.Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user != nil {
			t.Errorf("Failed to authenticate the user in the database : no error, whereas provided wrong email")
		}
	}
	err = handler.DeleteUser(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestAuthenticateWrongPassword(t *testing.T) {
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
	tmpUser, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	user, err := handler.Authenticate(tmpUser.Email, "Toto")
	if err != nil {
		t.Errorf("Failed to authenticate the user in the database : " + err.Error())
	} else {
		if user != nil {
			t.Errorf("Failed to authenticate the user in the database : no error, whereas provided wrong password")
		}
	}
	err = handler.DeleteUser(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}

func TestModifyUserRight(t *testing.T) {
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
	tmpUser, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		t.Errorf("Failed to add the user into the database : " + err.Error())
		t.FailNow()
	}
	user, err := handler.ModifyUser(tmpUser.ID, "toto@toto.com", "Toto", "Toto")
	if err != nil {
		t.Errorf("Failed to modify the user in the database : " + err.Error())
	} else {
		if user.Email != "toto@toto.com" {
			t.Errorf("Failed to modify the email of the user in the database")
		}
		if user.Name != "Toto" {
			t.Errorf("Failed to modify the name of the user in the database")
		}
		if user.Password != "Toto" {
			t.Errorf("Failed to modify the password of the user in the database")
		}
	}
	err = handler.DeleteUser(tmpUser.ID)
	if err != nil {
		t.Errorf("Failed to delete the user : " + err.Error())
	}
}
