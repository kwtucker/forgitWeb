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

//
//
// -=-=-=-=-=-= Covert times -=-=-=-=-=
// // Converts string Unix time stamp of UTC to int64
// dD, err := strconv.ParseInt(dbUser.LastUpdate, 10, 64)
// if err != nil {
// 	fmt.Println(err)
// }
// // conver to time.Time for compare
// dbUserUpdateDate := time.Unix(dD, 0)
// fmt.Println(dbUserUpdateDate)
// -=-=-=-=-=-= End Covert times -=-=-=-=-=

//
//
// =-=-=-=-API LOGIC =-=-=-=-

// if the local env update time is after the db update time, then take settings from local

// call update on user settings to db

// else send server user settings to the local env

// respond with the server settings

// =-=-=-=- END API LOGIC =-=-=-=-
