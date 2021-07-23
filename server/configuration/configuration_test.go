package configuration

import (
	"os"
	"testing"
)

func TestParseEnvironmentDatabaseURL(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	databaseURL, set := os.LookupEnv("DATABASE_URL")
	if !set {
		err := os.Setenv("DATABASE_URL", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.DatabaseURL != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("DATABASE_URL")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.DatabaseURL != databaseURL {
			t.Errorf("Failed to get the environment variable DATABASE_URL\n\tExpected : %q\n\tGiven : %q",
				databaseURL, conf.DatabaseURL)
		}
	}
}

func TestParseEnvironmentDatabaseUser(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	databaseUser, set := os.LookupEnv("DATABASE_USER")
	if !set {
		err := os.Setenv("DATABASE_USER", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.DatabaseUser != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("DATABASE_USER")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.DatabaseUser != databaseUser {
			t.Errorf("Failed to get the environment variable DATABASE_URL\n\tExpected : %q\n\tGiven : %q",
				databaseUser, conf.DatabaseUser)
		}
	}
}

func TestParseEnvironmentDatabasePassword(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	databasePassword, set := os.LookupEnv("DATABASE_PASSWORD")
	if !set {
		err := os.Setenv("DATABASE_PASSWORD", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.DatabasePassword != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("DATABASE_PASSWORD")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.DatabasePassword != databasePassword {
			t.Errorf("Failed to get the environment variable DATABASE_URL\n\tExpected : %q\n\tGiven : %q",
				databasePassword, conf.DatabasePassword)
		}
	}
}

func TestParseEnvironmentJWTSecretKey(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	jwtSecretKey, set := os.LookupEnv("JWT_SECRET_KEY")
	if !set {
		err := os.Setenv("JWT_SECRET_KEY", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.JWTSecretKey != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("JWT_SECRET_KEY")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.JWTSecretKey != jwtSecretKey {
			t.Errorf("Failed to get the environment variable JWT_SECRET_KEY\n\tExpected : %q\n\tGiven : %q",
				jwtSecretKey, conf.JWTSecretKey)
		}
	}
}

func TestParseEnvironmentDatabaseName(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	databaseName, set := os.LookupEnv("DATABASE_NAME")
	if !set {
		err := os.Setenv("DATABASE_NAME", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.DatabaseName != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("DATABASE_NAME")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.DatabaseName != databaseName {
			t.Errorf("Failed to get the environment variable DATABASE_URL\n\tExpected : %q\n\tGiven : %q",
				databaseName, conf.DatabaseName)
		}
	}
}

func TestParseEnvironmentTokenCookieName(t *testing.T) {
	conf, err := Parse("../configuration.json")
	if err != nil {
		t.Errorf("Failed to find the configuration.json file")
	}
	tokenCookieName, set := os.LookupEnv("TOKEN_COOKIE_NAME")
	if !set {
		err := os.Setenv("TOKEN_COOKIE_NAME", "toto")
		if err != nil {
			t.Errorf("Failed to set environment variable DATABASE_URL")
		}
		conf, err = Parse("../configuration.json")
		if err != nil {
			t.Errorf("Failed to find the configuration.json file")
		}
		if conf.TokenCookieName != "toto" {
			t.Errorf("Failed to get the environment variable DATABASE_URL")
		}
		err = os.Unsetenv("TOKEN_COOKIE_NAME")
		if err != nil {
			t.Errorf("Failed to unset the environment variable DATABASE_URL")
		}
	} else {
		if conf.TokenCookieName != tokenCookieName {
			t.Errorf("Failed to get the environment variable JWT_SECRET_KEY\n\tExpected : %q\n\tGiven : %q",
				tokenCookieName, conf.TokenCookieName)
		}
	}
}

func TestWrongPath(t *testing.T) {
	_, err := Parse("toto.json")
	if err == nil {
		t.Errorf("This path is wrong, no error has been returned, you need to fix it")
	}
}
