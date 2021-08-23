package main

import (
	"fmt"
	"teissem/stormtask/server/configuration"
)

func main() {
	conf, err := configuration.Parse("./configuration.json")
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print(conf)
	}
}
