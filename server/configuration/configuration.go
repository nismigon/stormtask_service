package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type ConfStruct struct {
	Address          string `json:"address"`
	Port             int    `json:"port"`
	DatabaseURL      string `json:"database_url"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	DatabaseName     string `json:"database_name"`
	JWTSecretKey     string `json:"jwt_secret_key"`
	TokenCookieName  string `json:"token_cookie_name"`
	BcryptCost       int    `json:"bcrypt_cost"`
}

// Parse analyze a configuration file and return the corresponding struct
// path	: Path to the configuration file
// This function return a pointer to the confirguration file
// This function can also return an error if it is not possible to open the file, to read the file
// or to unmarshal the JSON content
func Parse(path string) (*ConfStruct, error) {
	jsonFile, err := os.Open(filepath.Clean(path))
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
	conf.Address = setStringEnvVariable(conf.Address, "ADDRESS")
	portStr, set := os.LookupEnv("PORT")
	if set {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		conf.Port = port
	}
	conf.DatabaseURL = setStringEnvVariable(conf.DatabaseURL, "DATABASE_URL")
	conf.DatabaseUser = setStringEnvVariable(conf.DatabaseUser, "DATABASE_USER")
	conf.DatabasePassword = setStringEnvVariable(conf.DatabasePassword, "DATABASE_PASSWORD")
	conf.DatabaseName = setStringEnvVariable(conf.DatabaseName, "DATABASE_NAME")
	conf.JWTSecretKey = setStringEnvVariable(conf.JWTSecretKey, "JWT_SECRET_KEY")
	conf.TokenCookieName = setStringEnvVariable(conf.TokenCookieName, "TOKEN_COOKIE_NAME")
	bcryptCost, set := os.LookupEnv("BCRYPT_COST")
	if set {
		cost, err := strconv.Atoi(bcryptCost)
		if err != nil {
			return nil, err
		}
		conf.BcryptCost = cost
	}
	return &conf, nil
}

func setStringEnvVariable(current, environment string) string {
	variable, set := os.LookupEnv(environment)
	if set {
		return variable
	}
	return current
}
