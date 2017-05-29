package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/newestuser/contract-committee/app/rest"
	"github.com/newestuser/contract-committee/app"
	"github.com/newestuser/contract-committee/app/persistence/datastore/mongo"
	"fmt"
)

func main() {

	conf := mongo.Config{Hosts:[]string{"localhost"}, DBName:"committee", Indexes:[]*mongo.Index{}}
	db := mongo.NewDatabase(conf)

	suite := app.NewSuiteInteractor(db)

	r := mux.NewRouter()

	r.Handle("/tests", rest.RegisterTest(suite)).Methods("POST")
	r.Handle("/tests/{id}", nil).Methods("GET")
	r.Handle("/tests/{id}/cases", nil).Methods("POST")
	r.Handle("/tests/{testId}/cases/{caseId}/execute", nil).Methods("GET")

	go fmt.Println("contract-committee started on 8080")

	http.ListenAndServe(":8080", r)

}
