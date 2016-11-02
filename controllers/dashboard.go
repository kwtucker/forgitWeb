package controllers

import (
	"fmt"
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
		nerr, ncom, npush, rval int
		settingRepos            []models.SettingRepo
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

	apc, err := strconv.Atoi(r.Form.Get("apcMin"))
	if err != nil {
		log.Println(err)
	}
	p, err := strconv.Atoi(r.Form.Get("pMin"))
	if err != nil {
		log.Println(err)
	}

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

	var currentSetIndex int

	for i := range dbUser.Settings {
		if dbUser.Settings[i].Status == 1 {
			currentSetIndex = i
			break
		}
	}

	set := models.Setting{
		Name:   r.Form["settingGroupName"][0],
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
	setExists, err := c.db.SettingUserExistsCheck(dbconnect, session.Values["userID"].(int), dbUser.Settings[currentSetIndex].Name)
	if err != nil {
		log.Println(err)
	}

	dbUser.LastUpdate = "1"
	// add setting
	if setExists {
		// Add setting group to user settings
		dbUser.Settings[currentSetIndex] = set
		// Update user in db
		c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
	} else {
		// Add setting group to user settings
		dbUser.Settings = append(dbUser.Settings, set)
		// Update user in db
		// c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
		for i := range dbUser.Settings {
			if dbUser.Settings[i].Name == r.Form["settingGroupName"][0] {
				dbUser.Settings[i] = set
				c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
				break
			}
		}
	}

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
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

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
}

// SettingNew when the user pushes the new button on the dashboard form
// Creates a new setting group
func (c *DashboardController) SettingNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Copy db pipeline
	dbconnect := c.ConnectMongoDBStream()
	defer dbconnect.DBSession.Close()

	var (
		settingRepos []models.SettingRepo
	)

	// Grab most current user info
	dbUser, err := c.db.FindOneUser(dbconnect, session.Values["userID"].(int))
	if err != nil {
		log.Println(err)
	}

	for _, v := range dbUser.Repos {
		setrep := models.SettingRepo{
			GithubRepoID: v.RepoID,
			Name:         v.Name,
			Status:       0,
		}
		settingRepos = append(settingRepos, setrep)
	}

	set := models.Setting{
		Name:   "New (Change me)",
		Status: 1,
		SettingNotifications: models.SettingNotifications{
			OnError:  1,
			OnCommit: 1,
			OnPush:   1,
		},
		SettingAddPullCommit: models.SettingAddPullCommit{
			TimeMin: 5,
		},
		SettingPush: models.SettingPush{
			TimeMin: 60,
		},
		Repos: settingRepos,
	}
	setExists, err := c.db.SettingUserExistsCheck(dbconnect, session.Values["userID"].(int), "New (Change me)")
	if err != nil {
		log.Println(err)
	}

	// add setting
	if setExists == false {
		// Add setting group to user settings
		dbUser.Settings = append(dbUser.Settings, set)
		// Update user in db
		c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)

		for v := range dbUser.Settings {
			if dbUser.Settings[v].Status == 1 {
				dbUser.Settings[v].Status = 0
				c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
			}
			if dbUser.Settings[v].Name == set.Name {
				dbUser.Settings[v].Status = 1
				c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
			}
		}
	}

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
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
	fmt.Println(len(dbUser.Settings))
	if len(dbUser.Settings) > 1 {
		for i, v := range dbUser.Settings {
			if v.Status == 1 {
				dbUser.Settings = append(dbUser.Settings[:i], dbUser.Settings[i+1:]...)
				dbUser.Settings[0].Status = 1
				c.db.UpdateOneUser(dbconnect, session.Values["userID"].(int), &dbUser)
			}
		}
	}

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
}