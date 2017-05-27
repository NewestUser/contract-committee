package app

import (
	"github.com/newestuser/contract-committee/app/datastore"
)

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
	db datastore.DB
}

func (c *persistentCommittee) RegisterTest(t *NewTest) *Test {
	return &Test{}
}

func NewCommittee(db datastore.DB) Committee {
	return &persistentCommittee{db:db}
}
