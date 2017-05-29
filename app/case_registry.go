package app

import (
	"github.com/newestuser/contract-committee/app/assert"
)

type CaseRegistry interface {
	RegisterCase(suiteID string, c *assert.NewCase) (*assert.Case, error)
}

func NewCaseRegistry(v assert.Validator, s CaseStorage) CaseRegistry {
	return &persistentCaseRegistry{validator:v, storage:s}
}

type persistentCaseRegistry struct {
	storage   CaseStorage
	validator assert.Validator
}

func (r *persistentCaseRegistry) RegisterCase(suiteID string, c *assert.NewCase) (*assert.Case, error) {

	if err := r.validator.Valid(c); err != nil {
		return nil, BadReqErr("Cannot create invalid case: %s", err.Error())
	}

	id := r.storage.Register(suiteID, c)

	return &assert.Case{ID:id, NewCase:c}, nil
}