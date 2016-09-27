package controllers

import (
	// "fmt"
	// "github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	// "github.com/kwtucker/forgit/lib"
	// "github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system" // "log"
	// "net/http"
	// "time"
)

// APIController ...
type APIController struct {
	Env         system.Application
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// API ...
// func (c *APIController) API(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
// }
