package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	// "github.com/kwtucker/forgit/lib"
	"encoding/json"
	"github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	"log"
	"net/http"
	"strconv"
	// "time"
)

// APIController ...
type APIController struct {
	Env         system.Application
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// Connect will make a new copy of the main mongodb connection.
func (c *APIController) Connect() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy(), DName: c.DataConnect.DName}
}

//API ...
func (c *APIController) API(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dbconnect := c.Connect()
	defer dbconnect.DBSession.Close()

	var (
		err      error
		dbUser   models.User
		response []byte
		settings []models.APISetting
		resp     []models.APIError
	)

	ghid, err := strconv.Atoi(ps.ByName("ghid"))
	if err != nil {
		log.Println(err)
	}

	CheckUserExists, err := c.db.Exists(dbconnect, &ghid)
	if err != nil {
		log.Println("User was", err, "in the database - CheckUserExists")
	}

	switch CheckUserExists {
	case true:
		// Grab most current user info
		dbUser, err = c.db.FindOneUser(dbconnect, ghid)
		if err != nil {
			log.Println(err)
		}

		if dbUser.ForgitID == ps.ByName("fid") {

			for _, s := range dbUser.Settings {
				set := models.APISetting{
					Name:                 s.Name,
					Status:               s.Status,
					SettingNotifications: s.SettingNotifications,
					SettingAddPullCommit: s.SettingAddPullCommit,
					SettingPush:          s.SettingPush,
					Repos:                s.Repos,
				}
				json.Marshal(set)
				settings = append(settings, set)
			}
			response, err = json.Marshal(settings)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			res := models.APIError{
				Message: "bad credentials",
				Status:  http.StatusUnauthorized,
			}
			json.Marshal(res)
			resp = append(resp, res)
			response, err = json.Marshal(resp)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}

	case false:

		res := models.APIError{
			Message: "bad credentials",
			Status:  http.StatusUnauthorized,
		}
		json.Marshal(res)
		resp = append(resp, res)
		response, err = json.Marshal(resp)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

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
