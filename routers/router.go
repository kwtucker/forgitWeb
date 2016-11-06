package routers

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgitWeb/controllers"
	"github.com/kwtucker/forgitWeb/db"
	"github.com/kwtucker/forgitWeb/lib"
	"github.com/kwtucker/forgitWeb/system"
	"net/http"
)

// Init ...
func Init(application system.Application, router *httprouter.Router, database *db.ConnectMongo) {
	var sessionSecret = application.Config.SessionSecret
	var Store = sessions.NewCookieStore([]byte(sessionSecret))

	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir(application.Config.StaticPath))

	// -=--=-=- Routes -=-=-=-=-=
	// Root route for landing page
	router.GET("/", application.Route(
		&controllers.IndexController{
			Env:         application,
			Sess:        Store,
			DataConnect: database,
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

	// Dashboard Route that is what the callback goes too.
	router.GET("/dashboard/", application.Route(
		&controllers.DashboardController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "Dashboard"))

	//
	router.POST("/dashboard/setValues/", application.NoViewRoute(
		&controllers.DashboardController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "SettingSubmit"))

	router.POST("/dashboard/setSelect/", application.NoViewRoute(
		&controllers.DashboardController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "SettingSelect"))

	router.GET("/dashboard/setRemove/", application.NoViewRoute(
		&controllers.DashboardController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "SettingRemove"))

	router.GET("/getting-started/", application.Route(
		&controllers.GettingStartedController{
			Sess:        Store,
			Env:         application,
			DataConnect: database,
		}, "GettingStarted"))

	// API get user.
	router.GET("/api/users/:fid/:i", application.NoViewRoute(
		&controllers.APIController{
			Env:         application,
			DataConnect: database,
		}, "API"))

	// Logout will clear the sessions storage
	router.GET("/logout", application.NoViewRoute(
		&lib.Auth{
			Sess: Store,
			Env:  application,
		}, "Logout"))
}
