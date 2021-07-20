package web

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func TestInitAndCloseServer(t *testing.T) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to parse the configuration file : " + err.Error())
	}
	server, err := InitServer(*conf)
	if err != nil {
		t.Errorf("Failed to init the server : " + err.Error())
	}
	err = server.Close()
	if err != nil {
		t.Errorf("Failed to close the server : " + err.Error())
	}
}
