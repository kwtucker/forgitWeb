package controllers

import (
	// "encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/lib"
	// "github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	"golang.org/x/oauth2"
	// "io/ioutil"
	// "log"
	"net/http"
)

// TerminalController ...
type TerminalController struct {
	Env         system.Application
	Sess        *sessions.CookieStore
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// Connect will make a new copy of the main mongodb connection.
func (c *TerminalController) Connect() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy(), DName: c.DataConnect.DName}
}

// Terminal ...
func (c *TerminalController) Terminal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// This will check if user is authed
	var isAuth = &lib.Auth{Sess: c.Sess, Env: c.Env}

	// validate the session and get the session back
	session := isAuth.SessionCheck(w, r)

	// get access token and get client.
	tok, _ := c.Env.AuthConf.Exchange(oauth2.NoContext, session.Values["code"].(string))
	tokc := c.Env.AuthConf.Client(oauth2.NoContext, tok)
	client := github.NewClient(tokc)
	fmt.Println(client.Octocat("Oh JUMMMMMM"))

	// copy db pipeline and
	// don't close session tell end of function
	dbconnect := c.Connect()
	defer dbconnect.DBSession.Close()

	c.db.AddUser(dbconnect, client)

	// Nav for this view.
	navLinks := map[string]string{
		"/":               "Home",
		"/terminal/":      "Terminal",
		"#gettingstarted": "Getting Started",
		"/logout":         "Logout",
	}
	// data for the view
	data := map[string]interface{}{
		"PageName":        "Terminal",
		"ContentTemplate": "terminal",
		"NavLinks":        navLinks,
	}
	return data, http.StatusOK
}
