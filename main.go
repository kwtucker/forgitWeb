package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgitWeb/db"
	"github.com/kwtucker/forgitWeb/routers"
	"github.com/kwtucker/forgitWeb/system"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	// Defining which config file to parse though
	filename := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Setting the new application instance
	var application = &system.Application{}

	fmt.Printf("Using config %s\n", *filename)

	// Initialize app data
	application.Init(filename)

	// data is a new session to mongo and a database
	var database = &db.ConnectMongo{}
	database.ConnectDB(application.Config.DbHostString(), application.Config.DbName)

	// Create router server
	router := httprouter.New()

	// Call router function and send app and router to it.
	// Exporting router to router directory
	routers.Init(*application, router, database)

	// If route method is not allowed it will be a status erro
	router.HandleMethodNotAllowed = false

	// middleware
	handler := cors.Default().Handler(context.ClearHandler(router))
	// Starting the server with with port from the config file
	log.Printf("Listening on: %s", application.Config.HostString())
	log.Fatal(http.ListenAndServe(application.Config.WebPortString(), handler))
}
