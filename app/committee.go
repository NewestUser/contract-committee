package app

import "github.com/clouway/cloudplatform/task-queue/datastore/mongo"

type NewTest struct {
	Name string
}

type Test struct {
	Name string
	ID   string
}

type Committee interface {
	RegisterTest(t *NewTest) *Test
}

type persistentCommittee struct {
	db *mongo.Database
}

func (c *persistentCommittee) RegisterTest(t *NewTest) *Test {
	return &Test{}
}

func NewCommittee(db *mongo.Database) Committee {
	return &persistentCommittee{db:db}
}
