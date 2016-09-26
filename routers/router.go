package routers

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/controllers"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/lib"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// Init ...
func Init(application system.Application, router *httprouter.Router, database *db.ConnectMongo) {
	var sessionSecret = application.Config.SessionSecret
	var Store = sessions.NewCookieStore([]byte(sessionSecret))

	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// -=--=-=- Routes -=-=-=-=-=
	// Root route for landing page
	router.GET("/", application.Route(
		&controllers.IndexController{
			Env:  application,
			Sess: Store,
		}, "Index"))

	// Redirects to github for auth
	router.GET("/auth/", application.NoViewRoute(
		&lib.Auth{
			Env: application,
		}, "AuthFunc"))

	// Callback is what github goes to after auth
	router.GET("/auth/callback/", application.NoViewRoute(
		&lib.Auth{
			Env:  application,
			Sess: Store,
		}, "Callback"))

	// Terminal Route that is what the callback goes too.
	router.GET("/terminal/", application.Route(
		&controllers.TerminalController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "Terminal"))

	// Terminal Route that is what the callback goes too.
	router.GET("/getting-started/", application.Route(
		&controllers.GettingStartedController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "GettingStarted"))

	// Logout will clear the sessions storage
	router.GET("/logout", application.NoViewRoute(
		&lib.Auth{
			Sess: Store,
			Env:  application,
		}, "Logout"))
}
