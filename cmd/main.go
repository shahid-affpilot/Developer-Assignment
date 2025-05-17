package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/routes"
)

var cnf *config.Config

func init() {
	var err error
	cnf, err = config.LoadConfig() // get .env variables
	if err != nil {
		log.Println("Config func does not working well")
	}

	// connect to database
	database.ConnDB(cnf.Database)

	// add system admin if does not exist
	database.InitAdminUser(cnf.Admin)
}

func main() {
	fmt.Println("AffPilot Auth Service starting...")
	log.Println("Server initialized")

	r := mux.NewRouter()

	routes.SetupRoutes(r)

	// Use the port from config
	serverAddr := fmt.Sprintf(":%d", cnf.Server.Port)
	log.Printf("Server started at %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
