package controllers

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/db"
	"github.com/kwtucker/forgit/lib"
	// "github.com/kwtucker/forgit/models"
	"github.com/kwtucker/forgit/system"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	// "strconv"
	// "time"
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

// Terminal Controller
func (c *TerminalController) Terminal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {

	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// If the person is not authed send them back to home
	switch session.Values["authed"] {
	case 0, nil:
		session.Values["authed"] = 0
		session.Values["token"] = nil
		session.Save(r, w)
		return nil, http.StatusFound
	}

	// Copy db pipeline and
	// Defer close session tell end of function
	dbconnect := c.Connect()
	defer dbconnect.DBSession.Close()

	// Get token pointer and grab client to make requests
	tokpointer := lib.GetTokenStruct(session.Values["token"].(string))
	tokc := c.Env.AuthConf.Client(oauth2.NoContext, tokpointer)
	client := github.NewClient(tokc)

	// fmt.Println(client.Octocat("Oh JUMMMMMM"))

	// Get logged in user from github
	//Current User, *Response/ API call count, error
	ghuser, _, err := client.Users.Get("")
	if err != nil {
		fmt.Println("yep")
		log.Println(err)
	}

	// Get logged in users repos from github
	//[]*Repository, *Response, error
	repos, _, err := client.Repositories.List("", nil)
	if err != nil {
		log.Println(err)
	}

	// Set session value of UID
	session.Values["userID"] = *ghuser.ID
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Didn't Save userID")
		fmt.Println(err)
	}

	CheckUserExists, err := c.db.Exists(dbconnect, ghuser.ID)
	if err != nil {
		log.Println("User was", err, "in the database - CheckUserExists")
	}

	// If user doesn't exist create them
	switch CheckUserExists {
	case false:
		user := lib.CreateUser(ghuser, repos, nil)
		err = c.db.AddUser(dbconnect, user)
		if err != nil {
			log.Println(err)
		}
	}

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		fmt.Println(err)
	}

	// data for the view
	data := map[string]interface{}{
		"Auth":            true,
		"PageName":        "Terminal",
		"ContentTemplate": "terminal",
		"User":            dbUser,
	}
	return data, http.StatusOK
}

// func (c *TerminalController) SettingSubmit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
//
// 	// TODO:Grab from form value once terminal is there to set these values
// 	var settings = []models.Setting{}
// 	set := models.Setting{
// 		SettingID: 1,
// 		Name:      "Work",
// 		Status:    0,
// 		SettingNotifications: models.SettingNotifications{
// 			Status:   1,
// 			OnError:  1,
// 			OnCommit: 1,
// 			OnPush:   1,
// 		},
// 		SettingAddPullCommit: models.SettingAddPullCommit{
// 			Status:  1,
// 			TimeMin: 5,
// 		},
// 		SettingPush: models.SettingPush{
// 			Status:  1,
// 			TimeMin: 60,
// 		},
// 		// Repos: settingRepos,
// 	}
// 	settings = append(settings, set)
//
// 	// create user with structs ,update database with new data
// 	// lib.CreateUser(struct, struct, newsetting or nil)
// 	createUser := lib.CreateUser(ghuser, repos, settings)
// 	c.db.UpdateOne(dbconnect, session.Values["userID"].(int), createUser)
// 	fmt.Println("updated here you go")
//
// }
