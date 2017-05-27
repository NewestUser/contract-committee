package datastoretest

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"github.com/newestuser/contract-committee/app/datastore"
	"github.com/newestuser/contract-committee/app/datastore/mongo"
)

type DB struct {
	datastore.DB
	database *mgo.Database
	session  *mgo.Session
}

func NewDatabase() *DB {
	c := mongo.Config{Hosts: []string{"dev.telcong.com"}, DBName: "testDb", Indexes: []*mongo.Index{}}

	mgoDB := mongo.NewDatabase(c)

	s := connect(c.Hosts)
	database := s.DB(c.DBName)

	db := &DB{mgoDB, database, s}

	return db
}

// Closes db session
func (db *DB) Close() {
	db.database.Session.Close() //.session.Close()
}

// Clean is a helper function which is used in tests for cleaning up
// all existing collections in the database
func (db *DB) Clean() {
	col, err := db.database.CollectionNames()
	if err != nil {
		panic(err)
	}

	for _, v := range col {
		// skip dropping of system tables
		if !strings.Contains(v, "system.") {
			err := db.database.C(v).DropCollection()

			if err != nil {
				panic(err)
			}
		}
	}
}

func connect(hosts []string) *mgo.Session {
	info := &mgo.DialInfo{Addrs: hosts, Timeout: time.Second}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	return session
}
