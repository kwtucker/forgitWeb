package controllers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/system"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// IndexController ...
type IndexController struct {
	sess        *sessions.CookieStore
	Env         system.Application
	DataConnect *db.ConnectMongo
}

// Connect will make a new copy of the main mongodb connection.
func (c *IndexController) Connect() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy()}
}

// Index ...
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	session, err := c.sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(session.Values["authed"])
	if session.Values["authed"] != 1 {
		session.Values["authed"] = 0
	}
	session.Save(r, w)

	// copy db pipeline and
	// don't close session tell end of function
	dbconnect := c.Connect()
	defer dbconnect.DBSession.Close()

	// select the db and collection to use
	d := dbconnect.DBSession.DB("test").C("people")

	// Insert and handle error if any
	err := d.Insert(&Person{"Kevin", "777777777"})
	if err != nil {
		log.Fatal(err)
	}

	// find one in db and set to struct
	result := Person{}
	err = d.Find(bson.M{"name": "Kevin"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
	fmt.Println("Phone:", result.Name)

	// Nav for this view.
	navLinks := map[string]string{
		"#":          "Features",
		"#pricing":   "Pricing",
		"#createdby": "Created By",
		"/auth/":     "Log In / Register",
	}
	// values for the view.
	data := map[string]interface{}{
		"PageName":        "Forgit",
		"ContentTemplate": "index",
		"NavLinks":        navLinks,
	}
	return data, http.StatusOK
}
