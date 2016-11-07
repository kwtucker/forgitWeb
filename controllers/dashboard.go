package controllers

import (
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgitWeb/db"
	"github.com/kwtucker/forgitWeb/lib"
	"github.com/kwtucker/forgitWeb/models"
	"github.com/kwtucker/forgitWeb/system"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// DashboardController ...
type DashboardController struct {
	Env         system.Application
	Sess        *sessions.CookieStore
	DataConnect *db.ConnectMongo
	db          db.ConnectMongo
}

// ConnectMongoDBStream will make a new copy of the main mongodb connection.
func (c *DashboardController) ConnectMongoDBStream() *db.ConnectMongo {
	return &db.ConnectMongo{DBSession: c.DataConnect.DBSession.Copy(), DName: c.DataConnect.DName}
}

// Dashboard Controller
func (c *DashboardController) Dashboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {

	var (
		template string
		pageName string
	)
	pageName = "Dashboard"
	template = "dashboard"

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
	dbconnect := c.ConnectMongoDBStream()
	defer dbconnect.DBSession.Close()

	// Get token pointer and grab client to make requests
	tokpointer := lib.GetTokenStruct(session.Values["token"].(string))
	tokc := c.Env.AuthConf.Client(oauth2.NoContext, tokpointer)
	client := github.NewClient(tokc)

	// Get logged in user from github
	//Current User, *Response/ API call count, error
	ghuser, _, err := client.Users.Get("")
	if err != nil {
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
		log.Println(err)
	}

	CheckUserExists, err := c.db.UserExistsCheck(dbconnect, ghuser.ID)
	if err != nil {
		log.Println("User was", err, "in the database - CheckUserExists")
	}

	// If user doesn't exist create them.
	switch CheckUserExists {
	case false:
		user := lib.CreateUser(ghuser, repos, nil)
		err = c.db.AddUser(dbconnect, user)
		if err != nil {
			log.Println(err)
		}
		pageName = "Getting Started"
		template = "gettingStarted"
	}

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

	// data for the view
	data := map[string]interface{}{
		"Auth":            true,
		"PageName":        pageName,
		"ContentTemplate": template,
		"User":            dbUser,
	}
	return data, http.StatusOK
}

// SettingSubmit ...
func (c *DashboardController) SettingSubmit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Copy db pipeline and
	dbconnect := c.ConnectMongoDBStream()
	defer dbconnect.DBSession.Close()

	var (
		nerr, ncom, npush, rval    int
		settingRepos               []models.SettingRepo
		setExists, setExistsOnEdit bool
	)

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	for _, v := range dbUser.Repos {
		if r.Form.Get(*v.Name) != "" {
			rval = 1
		} else {
			rval = 0
		}

		setrep := models.SettingRepo{
			GithubRepoID: v.RepoID,
			Name:         v.Name,
			Status:       rval,
		}
		settingRepos = append(settingRepos, setrep)
	}

	apc, _ := strconv.Atoi(r.Form.Get("apcMin"))
	p, _ := strconv.Atoi(r.Form.Get("pMin"))

	if r.Form.Get("notifyErrors") != "" {
		nerr = 1
	} else {
		nerr = 0
	}

	if r.Form.Get("notifyCommit") != "" {
		ncom = 1
	} else {
		ncom = 0
	}

	if r.Form.Get("notifyPush") != "" {
		npush = 1
	} else {
		npush = 0
	}

	// Model out setting group with values
	set := models.Setting{
		Name:   strings.ToLower(r.Form["settingGroupName"][0]),
		Status: 1,
		SettingNotifications: models.SettingNotifications{
			OnError:  nerr,
			OnCommit: ncom,
			OnPush:   npush,
		},
		SettingAddPullCommit: models.SettingAddPullCommit{
			TimeMin: apc,
		},
		SettingPush: models.SettingPush{
			TimeMin: p,
		},
		Repos: settingRepos,
	}
	// Check if group exist for either update or new setting group, but doesn't check if new updated name.
	setExists, err = c.db.SettingUserExistsCheck(dbconnect, session.Values["userID"].(int), strings.ToLower(r.Form["settingGroupName"][0]))
	if err != nil {
		log.Println(err)
	}

	// If the setNameHide is something and not empty check if that name exists.
	if len(r.Form["setNameHide"]) != 0 {
		// If setting group is being updated and the name is changed, checks to see if previous name exitst
		setExistsOnEdit, err = c.db.SettingUserExistsCheck(dbconnect, session.Values["userID"].(int), strings.ToLower(r.Form["setNameHide"][0]))
		if err != nil {
			log.Println(err)
		}
	}

	dbUser.LastUpdate = "1"

	if setExists || setExistsOnEdit {
		// Add setting group to user settings
		// If settings name equal on in the DB
		for i := range dbUser.Settings {
			if dbUser.Settings[i].Name == strings.ToLower(r.Form["settingGroupName"][0]) || dbUser.Settings[i].Name == strings.ToLower(r.Form["setNameHide"][0]) {
				dbUser.Settings[i] = set
				break
			}
		}
		// Update user in db
		c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
	} else {
		// Add setting group to user settings
		dbUser.Settings = append(dbUser.Settings, set)

		// remove the current setting status to 0
		for i := range dbUser.Settings {
			if dbUser.Settings[i].Status == 1 {
				dbUser.Settings[i].Status = 0
				break
			}
		}
		c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
	}

	http.Redirect(w, r, "/dashboard/?s=true#settingGroups", http.StatusFound)
}

// SettingSelect ...
func (c *DashboardController) SettingSelect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Copy db pipeline and
	dbconnect := c.ConnectMongoDBStream()
	defer dbconnect.DBSession.Close()

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	for v := range dbUser.Settings {
		if dbUser.Settings[v].Status == 1 {
			dbUser.Settings[v].Status = 0
			c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
		}
		if dbUser.Settings[v].Name == r.Form["settingGroupSelect"][0] {
			dbUser.Settings[v].Status = 1
			c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
		}
	}

	http.Redirect(w, r, "/dashboard/#settingGroups", http.StatusFound)
}

//SettingRemove ...
func (c *DashboardController) SettingRemove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Copy db pipeline
	dbconnect := c.ConnectMongoDBStream()
	defer dbconnect.DBSession.Close()

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

	if len(dbUser.Settings) > 1 {
		for i, v := range dbUser.Settings {
			if v.Status == 1 {
				dbUser.Settings = append(dbUser.Settings[:i], dbUser.Settings[i+1:]...)
				dbUser.Settings[0].Status = 1
				c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
			}
		}
	}

	http.Redirect(w, r, "/dashboard/?r=true#settingGroups", http.StatusFound)
}
