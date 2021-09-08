package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Needed MySQL driver
)

type DBHandler struct {
	Handler    *sql.DB // Handler of the database connection
	URL        string  // Url to connect to the database
	User       string  // User to connect to the database
	Password   string  // Password to connect to the database
	Name       string  // Name of the database
	BcryptCost int     // Bcrypt cost
}

// Init initializes the connection to the database
// databaseURL 		The URL of the database
// databaseUser 	The user to connect to the database
// databasePassword The password to connect to the database
// databaseName 	The name of the database
func Init(databaseURL, databaseUser, databasePassword, databaseName string, bcryptCost int) (*DBHandler, error) {
	dataSource := databaseUser + ":" + databasePassword + "@tcp(" + databaseURL + ")/"
	// Open the connection with the database
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// Create or get the database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName)
	if err != nil {
		return nil, err
	}
	// Get an handler for the database
	dataSource = dataSource + databaseName
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	handler := &DBHandler{
		Handler:    db,
		URL:        databaseURL,
		User:       databaseUser,
		Password:   databasePassword,
		Name:       databaseName,
		BcryptCost: bcryptCost,
	}
	err = handler.UserInit()
	if err != nil {
		return nil, err
	}
	err = handler.GroupInit()
	if err != nil {
		return nil, err
	}
	err = handler.TaskInit()
	if err != nil {
		return nil, err
	}
	return handler, nil
}

// Close close the connection with the database
func (db *DBHandler) Close() error {
	return db.Handler.Close()
}
