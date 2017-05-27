package app

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type fakeTestRepository struct {
	mock.Mock
}

func (f *fakeTestRepository) Register(t *NewTest) string {
	return f.Called(t).String(0)
}

func TestRegisterANewTest(t *testing.T) {
	repository := new(fakeTestRepository)

	nt := &NewTest{Name:"foo"}
	id := "id"

	repository.On("Register", nt).Return(id)

	got := NewCommittee(repository).RegisterTest(nt)
	want := &Test{Name:nt.Name, ID:id}

	assert.Equal(t, want, got)
	repository.AssertExpectations(t)
}

