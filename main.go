package main

import (
	"fmt"
	"net/http"
	"strconv"
	"teissem/stormtask/server/configuration"
	"teissem/stormtask/server/web"
)

func main() {
	conf, err := configuration.Parse("./configuration.json")
	if err != nil {
		fmt.Println("Failed to parse the configuration : " + err.Error())
	}
	server, err := web.InitServer(*conf)
	if err != nil {
		fmt.Println("Failed to initialize the web server : " + err.Error())
	}
	http.Handle("/", server.Router)
	completeAddress := conf.Address + ":" + strconv.Itoa(conf.Port)
	fmt.Println("Server listening on " + completeAddress)
	err = http.ListenAndServe(completeAddress, nil)
	if err != nil {
		fmt.Println("An error occurred : " + err.Error())
	}
	fmt.Println("Server has been closed...")
}
