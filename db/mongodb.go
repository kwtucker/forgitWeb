package db

import (
	// "github.com/kwtucker/forgit/models"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
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
