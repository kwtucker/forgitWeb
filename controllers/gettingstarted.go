package controllers

import (
	// "fmt"
	// "github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	// "github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	// "github.com/kwtucker/forgit/lib"
	// "github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	// "golang.org/x/oauth2"
	// "log"
	// "net/http"
	// "time"
)

// GettingStartedController ...
type GettingStartedController struct {
	Env         system.Application
	Sess        *sessions.CookieStore
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// GettingStarted ...
// func (c *GettingStartedController) GettingStarted(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
//
// }
