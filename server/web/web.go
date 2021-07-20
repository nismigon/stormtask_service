package web

import (
	"naleakan/stormtask/configuration"
	"naleakan/stormtask/database"

	"github.com/gorilla/mux"
)

type Server struct {
	Router   *mux.Router
	Database *database.DBHandler
}

// InitServer initializes the http server and the database
func InitServer(configuration configuration.ConfStruct) (*Server, error) {
	router := mux.NewRouter()
	db, err := database.Init(
		configuration.DatabaseURL,
		configuration.DatabaseUser,
		configuration.DatabasePassword,
		configuration.DatabaseName)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router:   router,
		Database: db,
	}
	server.InitRoutes()
	return server, nil
}

// Close close the server connection
func (server *Server) Close() error {
	return server.Database.Close()
}
