package db

import (
	"fmt"
	"github.com/kwtucker/forgit/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// AddUser adds a user to the forgit DB and users Collection
func (c *ConnectMongo) AddUser(dbCopy *ConnectMongo, newuser *models.User) error {
	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")
	// Insert and handle error if any
	err := d.Insert(newuser)
	return err
}

// Exists looks for the githubID of the user in the forgit DB and users Collection.
func (c *ConnectMongo) Exists(dbCopy *ConnectMongo, userID *int) (bool, error) {
	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	err := d.Find(bson.M{"githubID": userID}).One(&result)
	if err != nil {
		return false, err
	}
	return true, err
}

// ExistsFID looks in the db for the requests forgitid and returns booleen.
func (c *ConnectMongo) ExistsFID(dbCopy *ConnectMongo, userID string) (bool, error) {

	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	err := d.Find(bson.M{"forgitid": userID}).One(&result)

	if err != nil {
		return false, err
	}
	return true, err
}

// SettingExists finds the user with githubID and looks for the setting group.
func (c *ConnectMongo) SettingExists(dbCopy *ConnectMongo, userID int, setName string) (bool, error) {
	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	err := d.Find(bson.M{"githubID": userID}).One(&result)
	if err != nil {
		return false, err
	}
	var r = false
	for v := range result.Settings {
		if result.Settings[v].Name == setName {
			r = true
		}
	}
	return r, err
}

// FindOneUser finds one user using the githubID.
func (c *ConnectMongo) FindOneUser(dbCopy *ConnectMongo, userID int) (models.User, error) {
	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	// find one in db and set to struct
	err := d.Find(bson.M{"githubID": userID}).One(&result)
	return result, err
}

// UpdateOne ...
func (c *ConnectMongo) UpdateOne(dbCopy *ConnectMongo, id int, user *models.User) error {
	// Find the current user
	userfind, err := c.FindOneUser(dbCopy, id)
	if err != nil {
		fmt.Println(err)
	}
	// update user with new github infor
	err = dbCopy.DBSession.DB("forgit").C("users").Update(userfind, user)
	return err
}
