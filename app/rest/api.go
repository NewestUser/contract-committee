package rest

import (
	"github.com/newestuser/contract-committee/app"
	"net/http"
	"io"
	"encoding/json"
)

func RegisterTest(c app.Committee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := &newTestDTO{}
		readJSON(r.Body, dto)

		newTest := &app.NewTest{Name:dto.Name}
		createdTest := c.RegisterTest(newTest)

		respDTO := &testDTO{ID:createdTest.ID, Name:createdTest.Name}
		writeJSON(w, respDTO)

	}
}

func readJSON(reader io.Reader, object interface{}) {
	json.NewDecoder(reader).Decode(object)
}

func writeJSON(writer http.ResponseWriter, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(object)
}

type newTestDTO struct {
	Name string `json:"name"`
}

type testDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}