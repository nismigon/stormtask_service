package database

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func TestInitRight(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	_, err = Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		t.Errorf("Failed to init the database connection : " + err.Error())
	}
}

func TestInitWrongURL(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	_, err = Init("toto", conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database URL given")
	}
}

func TestInitWrongUser(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	_, err = Init(conf.DatabaseURL, "toto", conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database user given")
	}
}

func TestInitWrongPassword(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	_, err = Init(conf.DatabaseURL, conf.DatabaseUser, "password", conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database password given")
	}
}

func TestInitWrongName(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	_, err = Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, "toto")
	if err == nil {
		t.Errorf("This test should fail : Wrong database name given")
	}
}
