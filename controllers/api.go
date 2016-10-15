package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
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
	)

	CheckUserExists, _ := c.db.ExistsFID(dbconnect, ps.ByName("fid"))

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

			err = d.Find(bson.M{"forgitid": ps.ByName("fid")}).One(&userfind)
			if err != nil {
				log.Println(err)
			}
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
			response, err = json.Marshal(res)
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
		response, err = json.Marshal(res)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
