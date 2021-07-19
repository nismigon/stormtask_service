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
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		t.Errorf("Failed to init the database connection : " + err.Error())
	} else {
		err = handler.Close()
		if err != nil {
			t.Errorf("Failed to close the database : " + err.Error())
		}
	}
}

func TestInitWrongURL(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
		t.FailNow()
	}
	handler, err := Init("toto", conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database URL given")
		err = handler.Close()
		if err != nil {
			t.Errorf("Failed to close the database : " + err.Error())
		}
		t.FailNow()
	}
}

func TestInitWrongUser(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
		t.FailNow()
	}
	handler, err := Init(conf.DatabaseURL, "toto", conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database user given")
		err = handler.Close()
		if err != nil {
			t.Errorf("Failed to close the database : " + err.Error())
		}
		t.FailNow()
	}
}

func TestInitWrongPassword(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, "password", conf.DatabaseName)
	if err == nil {
		t.Errorf("This test should fail : Wrong database password given")
		err = handler.Close()
		if err != nil {
			t.Errorf("Failed to close the database : " + err.Error())
		}
		t.FailNow()
	}
}

func TestInitWrongName(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, "toto")
	if err == nil {
		t.Errorf("This test should fail : Wrong database name given")
		err = handler.Close()
		if err != nil {
			t.Errorf("Failed to close the database : " + err.Error())
		}
		t.FailNow()
	}
}
