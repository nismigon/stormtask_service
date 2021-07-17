package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfStruct struct {
	DatabaseUrl      string `json:"database_url"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	DatabaseName     string `json:"database_name"`
}

// Parse analyze a configuration file and return the corresponding struct
// path	: Path to the configuration file
// This function return a pointer to the confirguration file
// This function can also return an error if it is not possible to open the file, to read the file
// or to unmarshal the JSON content
func Parse(path string) (*ConfStruct, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var conf ConfStruct
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		return nil, err
	}
	databaseUrl, set := os.LookupEnv("DATABASE_URL")
	if set {
		conf.DatabaseUrl = databaseUrl
	}
	databaseUser, set := os.LookupEnv("DATABASE_USER")
	if set {
		conf.DatabaseUser = databaseUser
	}
	databasePassword, set := os.LookupEnv("DATABASE_PASSWORD")
	if set {
		conf.DatabasePassword = databasePassword
	}
	databaseName, set := os.LookupEnv("DATABASE_NAME")
	if set {
		conf.DatabaseName = databaseName
	}
	return &conf, nil
}
