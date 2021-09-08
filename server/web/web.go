package web

import (
	"teissem/stormtask/server/configuration"
	"teissem/stormtask/server/database"

	"github.com/gorilla/mux"
)

type Server struct {
	Router        *mux.Router
	Database      *database.DBHandler
	Configuration configuration.ConfStruct
}

// InitServer initializes the http server and the database
func InitServer(configuration configuration.ConfStruct) (*Server, error) {
	router := mux.NewRouter()
	db, err := database.Init(
		configuration.DatabaseURL,
		configuration.DatabaseUser,
		configuration.DatabasePassword,
		configuration.DatabaseName,
		configuration.BcryptCost)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router:        router,
		Database:      db,
		Configuration: configuration,
	}
	server.InitRoutes()
	return server, nil
}

// Close close the server connection
func (s *Server) Close() error {
	return s.Database.Close()
}
