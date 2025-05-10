package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
)

var cnf *config.Config

func init() {
	// connect to database
	var err error
	cnf, err = config.LoadConfig()
	if err != nil {
		log.Println("Config func does not working well")
	}

	database.ConnDB(cnf.Database)

	// add system admin if does not exist
	database.InitAdminUser(cnf.Admin)
}

func main() {
	fmt.Println("AffPilot Auth Service starting...")
	log.Println("Server initialized")

	r := mux.NewRouter()

	

	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
