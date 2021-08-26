package database

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func (suite *DatabaseTestSuite) TestInit() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		suite.T().Errorf("Failed to init the database connection : " + err.Error())
	} else {
		err = handler.Close()
		if err != nil {
			suite.T().Errorf("Failed to close the database : " + err.Error())
		}
	}
}

func (suite *DatabaseTestSuite) TestInitWrongURL() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
		suite.T().FailNow()
	}
	handler, err := Init("toto", conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		suite.T().Errorf("This test should fail : Wrong database URL given")
		err = handler.Close()
		if err != nil {
			suite.T().Errorf("Failed to close the database : " + err.Error())
		}
		suite.T().FailNow()
	}
}

func (suite *DatabaseTestSuite) TestInitWrongUser() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
		suite.T().FailNow()
	}
	handler, err := Init(conf.DatabaseURL, "toto", conf.DatabasePassword, conf.DatabaseName)
	if err == nil {
		suite.T().Errorf("This test should fail : Wrong database user given")
		err = handler.Close()
		if err != nil {
			suite.T().Errorf("Failed to close the database : " + err.Error())
		}
		suite.T().FailNow()
	}
}

func (suite *DatabaseTestSuite) TestInitWrongPassword() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, "password", conf.DatabaseName)
	if err == nil {
		suite.T().Errorf("This test should fail : Wrong database password given")
		err = handler.Close()
		if err != nil {
			suite.T().Errorf("Failed to close the database : " + err.Error())
		}
		suite.T().FailNow()
	}
}

func (suite *DatabaseTestSuite) TestInitWrongName() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, "toto")
	if err == nil {
		suite.T().Errorf("This test should fail : Wrong database name given")
		err = handler.Close()
		if err != nil {
			suite.T().Errorf("Failed to close the database : " + err.Error())
		}
		suite.T().FailNow()
	}
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
