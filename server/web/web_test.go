package web

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WebTestSuite struct {
	suite.Suite
}

func (suite *WebTestSuite) TestInitAndCloseServer() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		suite.T().Errorf("Failed to init the server : " + err.Error())
	}
	err = server.Close()
	if err != nil {
		suite.T().Errorf("Failed to close the server : " + err.Error())
	}
}

func TestWeb(t *testing.T) {
	suite.Run(t, new(WebTestSuite))
}
