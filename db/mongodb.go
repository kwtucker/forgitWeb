package db

import (
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

// AddUser ...
func (c *ConnectMongo) AddUser(dbCopy *ConnectMongo, newuser *models.User) error {
	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")
	// Insert and handle error if any
	err := d.Insert(newuser)
	return err
}

// Exists ...
func (c *ConnectMongo) Exists(dbCopy *ConnectMongo, userID *int) (bool, error) {
	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	err := d.Find(bson.M{"githubID": userID}).One(&result)
	if err != nil {
		return false, err
	}
	return true, err
}

// FindOneUser ..
func (c *ConnectMongo) FindOneUser(dbCopy *ConnectMongo, userID int) (models.User, error) {
	// select the db and collection to use
	d := dbCopy.DBSession.DB("forgit").C("users")
	result := models.User{}
	// find one in db and set to struct
	err := d.Find(bson.M{"githubID": userID}).One(&result)
	return result, err
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
