package app

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type fakeCaseStorage struct {
	mock.Mock
}

func (f *fakeCaseStorage) Register(suiteID string, c *NewCase) string {
	return f.Called(suiteID, c).String(0)
}

func TestRegisterNewCase(t *testing.T) {
	storage := new(fakeCaseStorage)

	nc := &NewCase{
		GivenReq:&Given{URL:"fake.com", Method:"POST", Body:"foo"},
		AssertResp:&Assertion{StatusCode:200, Body:"foo"},
	}

	suiteID := "suiteId"

	storage.On("Register", suiteID, nc).Return("caseId")

	got := NewCaseRegistry(storage).RegisterCase(suiteID, nc)
	want := &Case{ID:"caseId", NewCase:nc}

	assert.Equal(t, want, got)

	storage.AssertExpectations(t)
}
