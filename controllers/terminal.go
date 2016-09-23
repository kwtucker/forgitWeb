package controllers

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/lib"
	"github.com/kwtucker/forgit/system"
	"golang.org/x/oauth2"
	"log"
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
	// copy db pipeline and
	// don't close session tell end of function
	dbconnect := c.Connect()
	defer dbconnect.DBSession.Close()

	// This will check if user is authed
	var isAuth = &lib.Auth{Sess: c.Sess, Env: c.Env}

	// validate the session and get the session back
	// if request not valid redirected to "/""
	session := isAuth.SessionCheck(w, r)

	// If the user id is set the user is logged in
	if session.Values["userID"] == nil {
		tok, _ := c.Env.AuthConf.Exchange(oauth2.NoContext, session.Values["code"].(string))
		tokc := c.Env.AuthConf.Client(oauth2.NoContext, tok)
		client := github.NewClient(tokc)
		fmt.Println(client.Octocat("Oh JUMMMMMM"))

		// Get logged in user
		//[]*Repository, *Response, error
		ghuser, _, err := client.Users.Get("")
		if err != nil {
			log.Println(err)
		}
		session.Values["userID"] = *ghuser.ID
		session.Save(r, w)

		// Get logged in users repos
		//[]*Repository, *Response, error
		repos, _, err := client.Repositories.List("", nil)
		if err != nil {
			log.Println(err)
		}
		//
		CheckUserExists, err := c.db.Exists(dbconnect, ghuser.ID)
		if err != nil {
			log.Println(err)
		}

		switch CheckUserExists {
		case false:
			User := lib.CreateUser(ghuser, repos)
			err = c.db.AddUser(dbconnect, User)
			if err != nil {
				log.Println(err)
			}
		}
	}

	User, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

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
		"User":            User,
	}
	return data, http.StatusOK
}
