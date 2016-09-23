package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	config "github.com/SivaShhankar/CMS_NEW/Database"
	routers "github.com/SivaShhankar/CMS_NEW/Routers"
)

func main() {

	// Load the configuration stuffs.
	config.LoadAppConfig()

	// Initiate the database information.
	config.CreateDBSession()

	// Add the neccessary indexes.
	config.AddIndexes()

	// Created the routes of this application/
	mux := http.NewServeMux()
	mux = routers.SetCandidateRoutes(mux)

	log.Println("Listening...")

	// Listen the server.
	http.ListenAndServe(GetPort(), mux)
}

// GetPort -- get the Port from the Dynamic environment
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Println("Running Port is" + port)
	return ":" + port
}
