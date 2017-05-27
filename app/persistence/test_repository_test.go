package persistence

import (
	"testing"
	"os"
	"github.com/newestuser/contract-committee/app/persistence/datastore/datastoretest"
)

var db *datastoretest.DB

func TestMain(m *testing.M) {
	db = datastoretest.NewDatabase()

	code := m.Run()
	db.Close()
	os.Exit(code)
}
