package app

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type fakeSuiteRepository struct {
	mock.Mock
}

func (f *fakeSuiteRepository) Register(t *NewSuite) string {
	return f.Called(t).String(0)
}

func (f*fakeSuiteRepository) Exists(name string) bool {
	return f.Called(name).Bool(0)
}

func (f *fakeSuiteRepository) Find(id string) *Suite {
	if stub := f.Called(id).Get(0); stub != nil {
		return stub.(*Suite)
	}

	return nil
}

func TestRegisterANewTest(t *testing.T) {
	repository := new(fakeSuiteRepository)

	nt := &NewSuite{Name:"foo"}
	id := "id"

	repository.On("Exists", nt.Name).Return(false)
	repository.On("Register", nt).Return(id)

	got, err := NewSuiteInteractor(repository).RegisterSuite(nt)
	want := &Suite{Name:nt.Name, ID:id}

	assert.Nil(t, err)
	assert.Equal(t, want, got)

	repository.AssertExpectations(t)
}

func TestRegisterWithAlreadyExistingName(t *testing.T) {
	repo := new(fakeSuiteRepository)

	nt := &NewSuite{Name:"existing"}

	repo.On("Exists", nt.Name).Return(true)

	gotT, gotErr := NewSuiteInteractor(repo).RegisterSuite(nt)
	wantErr := BadReqErr("Test with name %s already exists", nt.Name)

	assert.Nil(t, gotT)
	assert.Equal(t, wantErr, gotErr)

	repo.AssertExpectations(t)
}

func TestFindExistingTest(t *testing.T) {
	repo := new(fakeSuiteRepository)

	et := &Suite{Name:"what ever", ID:"existingId"}

	repo.On("Find", et.ID).Return(et)

	gotTest, gotErr := NewSuiteInteractor(repo).FindSuite(et.ID)

	assert.Nil(t, gotErr)
	assert.Equal(t, et, gotTest)

	repo.AssertExpectations(t)
}

func TestDoNotFindTest(t *testing.T) {
	repo := new(fakeSuiteRepository)
	id := "missingId"

	repo.On("Find", id).Return(nil)

	gotTest, gotErr := NewSuiteInteractor(repo).FindSuite(id)
	wantErr := NotFoundErr("Test with id %s can not be found", id)

	assert.Nil(t, gotTest)
	assert.Equal(t, wantErr, gotErr)

	repo.AssertExpectations(t)
}