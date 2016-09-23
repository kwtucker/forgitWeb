package db

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/kwtucker/forgit/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// ConnectMongo ...
type ConnectMongo struct {
	DBSession *mgo.Session
	DName     *mgo.Database
}

// ConnectDB ...
func (c *ConnectMongo) ConnectDB(dburl, dbname string) {
	session, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	c.DBSession = session
	c.DName = session.DB(dbname)
}

// AddUser ...
func (c *ConnectMongo) AddUser(dbCopy *ConnectMongo, client *github.Client) {

	// Get logged in user
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Println(err)
	}

	repos, _, err := client.Repositories.List("", nil)
	if err != nil {
		fmt.Println(err)
	}

	var (
		repoArr      = []models.Repo{}
		settingRepos = []models.SettingRepo{}
		settings     = []models.Setting{}
	)

	for k, _ := range repos {
		currentUserRepos := models.Repo{
			URL:             repos[k].URL,
			CommitsURL:      repos[k].CommitsURL,
			ContributorsURL: repos[k].ContributorsURL,
			Description:     repos[k].Description,
			FullName:        repos[k].FullName,
			GitCommitsURL:   repos[k].GitCommitsURL,
			HTMLURL:         repos[k].HTMLURL,
			RepoID:          repos[k].ID,
			Name:            repos[k].Name,
			Owner:           repos[k].Owner.Login,
		}
		repoArr = append(repoArr, currentUserRepos)

		currentUserSettingsRepo := models.SettingRepo{
			GithubRepoID: repos[k].ID,
			Name:         repos[k].Name,
			Status:       0,
		}
		settingRepos = append(settingRepos, currentUserSettingsRepo)
	}

	currentUserSettings := models.Setting{
		SettingID: 1,
		Name:      "General",
		Status:    1,
		SettingNotifications: models.SettingNotifications{
			Status:   1,
			OnError:  1,
			OnCommit: 1,
			OnPush:   1,
		},
		SettingAddPullCommit: models.SettingAddPullCommit{
			Status:  1,
			TimeMin: 5,
		},
		SettingPush: models.SettingPush{
			Status:  1,
			TimeMin: 60,
		},
		Repos: settingRepos,
	}
	settings = append(settings, currentUserSettings)

	timenow := &github.Timestamp{time.Now()}
	currentUser := &models.User{
		GithubID:   user.ID,
		LastUpdate: timenow.String(),
		LastSync:   timenow.String(),
		Login:      user.Login,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		Company:    user.Company,
		HTMLURL:    user.HTMLURL,
		ReposURL:   user.ReposURL,
		Repos:      repoArr,
		Settings:   settings,
	}

	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")
	// find one in db and set to struct
	result := models.User{}
	// Insert and handle error if any
	err = d.Insert(currentUser)
	if err != nil {
		log.Fatal(err)
	}
	// err = d.Find(nil).Distinct("githubID", &result)
	err = d.Find(bson.M{"githubID": user.ID}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User Full Name:", *result.Name)
}

// Exists ...
func (c *ConnectMongo) Exists() {

}

// FindOne ..
func (c *ConnectMongo) FindOne() {

}

// FindAll ..
func (c *ConnectMongo) FindAll() {

}

// UpdateOne ..
func (c *ConnectMongo) UpdateOne() {

}

// Remove ..
func (c *ConnectMongo) Remove() {

}
