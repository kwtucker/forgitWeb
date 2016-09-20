package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/routers"
	"github.com/kwtucker/forgit/system"
	"log"
	"net/http"
)

func main() {
	filename := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()
	var application = &system.Application{}

	fmt.Printf("Using config %s\n", *filename)

	// Initialize app data
	application.Init(filename)

	// Create router server
	router := httprouter.New()

	// Call router function and send app and router to it.
	// Exporting router to router directory
	routers.Init(*application, router)

	/*
	 * If HandleMethodNotAllowed is enabled, the router checks if another method is allowed for the
	 * current route, if the current request can not be routed.
	 * If this is the case, the request is answered with 'Method Not Allowed'
	 * and HTTP status code 405.
	 * If no other Method is allowed, the request is delegated to the NotFound
	 * handler.
	 */
	router.HandleMethodNotAllowed = false

	// Starting the server with with port from the config file
	log.Printf("Listening on: %s", application.Config.HostString())
	log.Fatal(http.ListenAndServe(application.Config.WebPortString(), context.ClearHandler(router)))
}
