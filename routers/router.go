package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/controllers"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// Init ...
func Init(application system.Application, router *httprouter.Router, database *db.ConnectMongo) {
	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	// Routes
	router.GET("/", application.Route(&controllers.IndexController{Env: application, DataConnect: database}, "Index"))
}
