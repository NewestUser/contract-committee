package mongo

import (
	"gopkg.in/mgo.v2"
	"time"

	"github.com/clouway/cloudplatform/datastore"
)

// Config is the database configuration which is used for initialisation of the datastore.
type Config struct {
	// List of hosts to be used for communication
	Hosts []string
	// Name if the database
	DBName string
	// List of indexes to be initialised on start up
	Indexes []*Index
}

// Index a single index in the datastore.
type Index struct {
	// Kind to which index to be applied
	Kind string
	// Single or Composite Key which to be used for indexing
	Key []string
}

// A DB accesses the datastore (MongoDB).
type mongoDB struct {
	session  *mgo.Session
	database *mgo.Database
}

// New creates a new datastore by using the provided configuration.
func NewDatabase(c Config) datastore.DB {
	s := connect(c.Hosts)
	d := s.DB(c.DBName)

	db := &mongoDB{database: d, session: s}

	err := db.buildIndexes(c.Indexes)

	if err != nil {
		panic(err)
	}

	return db
}

func (d *mongoDB) C(name string) datastore.Collection {
	c := d.database.C(name)

	return &mongoCollection{session: d.session, collection: c}
}

func (d *mongoDB) buildIndexes(index []*Index) error {
	in := indexer{database: d.database}

	for _, cur := range index {
		in.createIndex(cur.Kind, mgo.Index{
			Key:        cur.Key,
			Unique:     false,
			DropDups:   false,
			Background: true,
			Sparse:     false})

	}

	return in.Error()
}

// Collection is wrapper of mgo.Collection, that automatically recycles used sessions
type mongoCollection struct {
	session    *mgo.Session    // the initial session, created from mgo.Dial
	collection *mgo.Collection // the initial collection, created from the initial session
}

func (mc *mongoCollection) Insert(doc interface{}) error {
	session, collection := mc.refresh()
	defer session.Close()

	return collection.Insert(doc)
}

func (mc *mongoCollection) Upsert(selector interface{}, update interface{}) error {
	session, collection := mc.refresh()
	defer session.Close()

	_, err := collection.Upsert(selector, update)

	return err
}

func (mc *mongoCollection) Aggregate(pipeline interface{}, result interface{}) error {
	session, collection := mc.refresh()
	defer session.Close()

	pipe := collection.Pipe(pipeline)
	iter := pipe.Iter()

	if iter.Err() != nil {
		return iter.Err()
	}

	return iter.All(result)
}

func (mc *mongoCollection) FindSorted(query interface{}, sort string, limit int, result interface{}) error {
	session, collection := mc.refresh()
	defer session.Close()

	return collection.Find(query).Limit(limit).Sort(sort).All(result)
}

func (mc *mongoCollection) Find(query interface{}, result interface{}) error {
	session, collection := mc.refresh()
	defer session.Close()

	return collection.Find(query).All(result)
}

func (mc *mongoCollection) refresh() (*mgo.Session, *mgo.Collection) {
	s := mc.session.Copy()
	c := mc.collection.With(s)

	return s, c
}

// indexer is an indexer which is used as helper for creation of indexes
type indexer struct {
	database *mgo.Database
	err      error
}

func (i *indexer) createIndex(coll string, index mgo.Index) {
	// some index before this thrown an error, so we have to skip this part
	if i.err != nil {
		return
	}

	err := i.database.C(coll).EnsureIndex(index)
	i.err = err
}

func (i indexer) Error() error {
	return i.err
}

func connect(hosts []string) *mgo.Session {
	info := &mgo.DialInfo{Addrs: hosts, Timeout: time.Second}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	return session
}
