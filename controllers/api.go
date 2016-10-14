package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	// "github.com/kwtucker/forgit/lib"
	"encoding/json"
	"github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	// "strconv"
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

	// ghid, err := strconv.Atoi(ps.ByName("ghid"))
	// if err != nil {
	// 	log.Println(err)
	// }

	// CheckUserExists, err := c.db.Exists(dbconnect, &ghid)
	CheckUserExists, err := c.db.ExistsFID(dbconnect, ps.ByName("fid"))

	if err != nil {
		log.Println("User was", err, "in the database - CheckUserExists")
	}

	switch CheckUserExists {
	case true:
		d := dbconnect.DBSession.DB("forgit").C("users")
		dbUser = models.User{}
		// find one in db and set to struct
		err = d.Find(bson.M{"forgitid": ps.ByName("fid")}).One(&dbUser)

		if dbUser.ForgitID == ps.ByName("fid") && dbUser.LastUpdate == "1" ||
			dbUser.ForgitID == ps.ByName("fid") && ps.ByName("i") == "init" {

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
			dbUser.LastUpdate = "0"
			d := dbconnect.DBSession.DB("forgit").C("users")
			userfind := models.User{}
			// err := d.Find(bson.M{"githubID": userID}).One(&result)
			err := d.Find(bson.M{"forgitid": ps.ByName("fid")}).One(&userfind)
			if err != nil {
				log.Println(err)
			}
			// c.db.UpdateOne(dbconnect, ghid, &dbUser)
			// update user with new github infor
			err = dbconnect.DBSession.DB("forgit").C("users").Update(userfind, dbUser)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else if dbUser.ForgitID == ps.ByName("fid") && dbUser.LastUpdate == "0" {
			upd := models.UpdateStatus{
				Update: "0",
			}
			response, err = json.Marshal(upd)
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
