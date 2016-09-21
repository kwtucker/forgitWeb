package routers

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/controllers"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// Init ...
func Init(application system.Application, router *httprouter.Router, database *db.ConnectMongo) {
	var sessionSecret = application.Config.SessionSecret
	var Store = sessions.NewCookieStore([]byte(sessionSecret))
	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// Routes
	router.GET("/", application.Route(&controllers.IndexController{Env: application, DataConnect: database, sess: Store}, "Index"))
	router.GET("/auth/", application.NoViewRoute(&lib.Auth{Env: application}, "AuthFunc"))
	router.GET("/terminal/", application.Route(&controllers.TerminalController{Env: application, sess: Store}, "Terminal"))
}
