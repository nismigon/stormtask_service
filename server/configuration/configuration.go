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
	address, set := os.LookupEnv("ADDRESS")
	if set {
		conf.Address = address
	}
	portStr, set := os.LookupEnv("PORT")
	if set {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		conf.Port = port
	}
	databaseURL, set := os.LookupEnv("DATABASE_URL")
	if set {
		conf.DatabaseURL = databaseURL
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
	jwtSecretKey, set := os.LookupEnv("JWT_SECRET_KEY")
	if set {
		conf.JWTSecretKey = jwtSecretKey
	}
	tokenCookieName, set := os.LookupEnv("TOKEN_COOKIE_NAME")
	if set {
		conf.TokenCookieName = tokenCookieName
	}
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
