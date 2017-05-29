package app

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"

	ass "github.com/newestuser/contract-committee/app/assert"
	"errors"
)

type fakeCaseStorage struct {
	mock.Mock
}

func (f *fakeCaseStorage) Register(suiteID string, c *ass.NewCase) string {
	return f.Called(suiteID, c).String(0)
}

type fakeCaseValidator struct {
	mock.Mock
}

func (f *fakeCaseValidator) Valid(c *ass.NewCase) error {
	return f.Called(c).Error(0)
}

func TestRegisterNewCase(t *testing.T) {
	storage := new(fakeCaseStorage)
	validator := new(fakeCaseValidator)

	nc := &ass.NewCase{
		GivenReq:&ass.Given{URL:"fake.com", Method:"POST", Body:"foo"},
		AssertResp:&ass.Assertion{StatusCode:200, Body:"foo"},
	}

	suiteID := "suiteId"

	validator.On("Valid", nc).Return(nil)
	storage.On("Register", suiteID, nc).Return("caseId")

	gotCase, gotErr := NewCaseRegistry(validator, storage).RegisterCase(suiteID, nc)
	wantCase := &ass.Case{ID:"caseId", NewCase:nc}

	assert.Nil(t, gotErr)
	assert.Equal(t, wantCase, gotCase)

	validator.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestRegisterInvalidCase(t *testing.T) {
	storage := new(fakeCaseStorage)
	validator := new(fakeCaseValidator)

	nc := &ass.NewCase{}

	validationError := errors.New("some error")

	validator.On("Valid", nc).Return(validationError)

	gotCase, gotErr := NewCaseRegistry(validator, storage).RegisterCase("anyId", nc)
	wantErr := BadReqErr("Cannot create invalid case: %s", validationError.Error())

	assert.Nil(t, gotCase)
	assert.Equal(t, wantErr, gotErr)

	validator.AssertExpectations(t)
	storage.AssertExpectations(t)
}