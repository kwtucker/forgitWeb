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

// IndexController ...
type IndexController struct {
	Sess        *sessions.CookieStore
	Env         system.Application
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// ConnectMongoDBStream will make a new copy of the main mongodb connection.
func (c *IndexController) ConnectMongoDBStream() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy(), DName: c.DataConnect.DName}
}

// Index landing view
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// Get Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// if the session is new make authed to 0
	if session.IsNew {
		session.Values["authed"] = 0
		session.Save(r, w)
	}
	var (
		authed bool
		dbUser models.User
	)
	if session.Values["authed"] == 1 {
		// Copy db pipeline and
		dbconnect := c.ConnectMongoDBStream()
		defer dbconnect.DBSession.Close()

		// Grab most current user info
		dbUser, err = c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
		if err != nil {
			log.Println(err)
		}
		authed = true
	} else {
		authed = false
	}

	// values for the view.
	data := map[string]interface{}{
		"Auth":            authed,
		"PageName":        "Forgit",
		"ContentTemplate": "index",
		"User":            dbUser,
	}
	return data, http.StatusOK
}
