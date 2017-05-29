package app

import "github.com/newestuser/contract-committee/app/assert"

type SuiteStorage interface {
	Register(t *NewSuite) string

	Exists(name string) bool

	Find(id string) *Suite
}

type CaseStorage interface {
	Register(suiteID string, c *assert.NewCase) string
}