package controllers

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgitWeb/db"
	"github.com/kwtucker/forgitWeb/models"
	"github.com/kwtucker/forgitWeb/system"
	"log"
	"net/http"
)

// DonateController  (Env, Sess, DataConnect, db)
type DonateController struct {
	Env         system.Application
	Sess        *sessions.CookieStore
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// ConnectMongoDBStream will make a new copy of the main mongodb connection.
func (c *DonateController) ConnectMongoDBStream() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy(), DName: c.DataConnect.DName}
}

// Donate ...
func (c *DonateController) Donate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var auth = true
	var dbUser models.User
	// If the person is not authed send them back to home
	switch session.Values["authed"] {
	case 0:
		auth = false
	}

	if auth == true {

		// Copy db pipeline and
		dbconnect := c.ConnectMongoDBStream()
		defer dbconnect.DBSession.Close()

		// Grab most current user info
		dbUser, err = c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
		if err != nil {
			log.Println(err)
		}
	}

	// data for the view
	data := map[string]interface{}{
		"Auth":            auth,
		"PageName":        "Donate",
		"ContentTemplate": "donate",
		"User":            dbUser,
	}
	return data, http.StatusOK
}
