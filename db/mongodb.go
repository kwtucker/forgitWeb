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

	var repoArr = []models.Repo{}
	for k, _ := range repos {
		currentUserRepos := models.Repo{
			URL:             repos[k].URL,
			CommitsURL:      repos[k].CommitsURL,
			ContributorsURL: repos[k].ContributorsURL,
			Description:     repos[k].Description,
			FullName:        repos[k].FullName,
			GitCommitsURL:   repos[k].GitCommitsURL,
			HTMLURL:         repos[k].HTMLURL,
			ID:              repos[k].ID,
			Name:            repos[k].Name,
			Owner:           repos[k].Owner.Login,
		}
		repoArr = append(repoArr, currentUserRepos)
	}
	timenow := &github.Timestamp{time.Now()}
	currentUser := &models.User{
		ID:         user.ID,
		LastUpdate: timenow.String(),
		LastSync:   timenow.String(),
		Login:      user.Login,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		Company:    user.Company,
		HTMLURL:    user.HTMLURL,
		ReposURL:   user.ReposURL,
		Repos:      repoArr,
		// Settings:   []Setting,
	}

	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")

	// Insert and handle error if any
	err = d.Insert(currentUser)
	if err != nil {
		log.Fatal(err)
	}

	// find one in db and set to struct
	result := models.User{}
	err = d.Find(bson.M{"login": "kwtucker"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User Full Name:", result.Name)

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

// Exists ...
func (c *ConnectMongo) Exists() {

}
