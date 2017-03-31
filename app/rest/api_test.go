package rest

import (
	"testing"
	"encoding/json"
	"bytes"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/mock"
	"github.com/newestuser/contract-committee/app"
	"io"
	"github.com/stretchr/testify/assert"
)

type fakeCommittee struct {
	mock.Mock
}

func (f *fakeCommittee) RegisterTest(t *app.NewTest) *app.Test {
	return f.Called(t).Get(0).(*app.Test)
}

func TestRegisterTest(t *testing.T) {
	dto := newTestDTO{Name: "my-test"}
	testReq := &app.NewTest{Name:dto.Name}
	testResp := &app.Test{ID:"id", Name:dto.Name}

	request := createJSONRequest(dto)
	recorder := httptest.NewRecorder()

	committee := new(fakeCommittee)
	committee.On("RegisterTest", testReq).Return(testResp)

	RegisterTest(committee).ServeHTTP(recorder, request)

	got := &testDTO{}
	readJSONResponse(recorder.Body, got)

	assert.Equal(t, dto.Name, got.Name)
	assert.Equal(t, testResp.ID, got.ID)

	committee.AssertExpectations(t)
}

func createJSONRequest(object interface{}) *http.Request {
	data, _ := json.Marshal(object)
	reader := bytes.NewReader(data)

	r, _ := http.NewRequest("POST", "http://contract-committee/", reader)

	return r
}

func readJSONResponse(reader io.Reader, object interface{}) {
	decoder := json.NewDecoder(reader)
	decoder.Decode(object)
}
