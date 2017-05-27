package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/newestuser/contract-committee/app/rest"
	"github.com/newestuser/contract-committee/app"
	"github.com/newestuser/contract-committee/app/datastore/mongo"
)

func main() {

	conf := mongo.Config{Hosts:[]string{"localhost"}, DBName:"committee", Indexes:[]*mongo.Index{}}
	db := mongo.NewDatabase(conf)

	committee := app.NewCommittee(db)

	r := mux.NewRouter()

	r.Handle("/tests", rest.RegisterTest(committee)).Methods("POST")

	http.ListenAndServe(":8080", r)
}
