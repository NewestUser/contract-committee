package datastoretest

import (
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"testing"
)

type Person struct {
	Name    string `bson:"name"`
	Address string `bson:"address"`
	Age     int    `bson:"age"`
}

var db *DB

func TestMain(m *testing.M) {
	db = NewDatabase()

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestInsert(t *testing.T) {
	db.Clean()

	p := &Person{Name: "Emil", Address: "Ivan Vazov", Age: 27}

	db.C("persons").Insert(p)

	got := make([]Person, 0)
	q := bson.M{}

	err := db.C("persons").Find(q, &got)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := []Person{Person{Name: "Emil", Address: "Ivan Vazov", Age: 27}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v", want)
		t.Errorf("     got %v", got)
	}
}

func TestUpsertAsInsert(t *testing.T) {
	db.Clean()

	p := &Person{Name: "Emil", Address: "Ivan Vazov", Age: 27}
	q := bson.M{"_id": "12345"}

	err := db.C("persons").Upsert(q, p)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	got := make([]Person, 0)
	err = db.C("persons").Find(bson.M{}, &got)

	if err != nil {
		t.Errorf("Unexapected error: %v", err)
	}

	want := []Person{Person{Name: "Emil", Address: "Ivan Vazov", Age: 27}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v", want)
		t.Errorf("     got %v", got)
	}
}

func TestUpsertAsUpdate(t *testing.T) {
	db.Clean()

	p := &Person{Name: "Emil", Address: "Ivan Vazov", Age: 26}
	db.C("persons").Insert(p)

	q := bson.M{"name": "Emil"}

	u := bson.M{"$set": bson.M{"age": 27}}
	err := db.C("persons").Upsert(q, u)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	got := make([]Person, 0)
	err = db.C("persons").Find(bson.M{}, &got)

	if err != nil {
		t.Errorf("Unexapected error: %v", err)
	}

	want := []Person{Person{Name: "Emil", Address: "Ivan Vazov", Age: 27}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v", want)
		t.Errorf("     got %v", got)
	}
}

func TestAggregate(t *testing.T) {
	db.Clean()

	db.C("persons").Insert(Person{Name: "Emil", Address: "Ivan Vazov", Age: 27})
	db.C("persons").Insert(Person{Name: "Petar", Address: "Ivan Vazov", Age: 26})
	db.C("persons").Insert(Person{Name: "Ivan", Address: "Vasil Levski", Age: 24})
	db.C("persons").Insert(Person{Name: "Marian", Address: "Hristo Botev", Age: 27})

	pipeline := []bson.M{
		{"$match": bson.M{"address": "Ivan Vazov"}},
		{"$group": bson.M{"_id": "$address"}},
	}

	got := []bson.M{}

	db.C("persons").Aggregate(pipeline, &got)

	want := "Ivan Vazov"

	if !reflect.DeepEqual(got[0]["_id"], want) {
		t.Errorf("expected %v", want)
		t.Errorf("     got %v", got[0]["_id"])
	}
}

func TestFindSorted(t *testing.T) {
	db.Clean()

	db.C("persons").Insert(Person{"Emil", "Ivan Vazov", 27})
	db.C("persons").Insert(Person{"Ivan", "Ivan Vazov", 26})
	db.C("persons").Insert(Person{"Petar", "Ivan Vazov", 23})
	db.C("persons").Insert(Person{"Marian", "Ivan Vazov", 24})

	result := make([]Person, 0)

	q := bson.M{"address": "Ivan Vazov"}

	db.C("persons").FindSorted(q, "age", 4, &result)

	expected := []Person{
		Person{"Petar", "Ivan Vazov", 23},
		Person{"Marian", "Ivan Vazov", 24},
		Person{"Ivan", "Ivan Vazov", 26},
		Person{"Emil", "Ivan Vazov", 27},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v", expected)
		t.Errorf("     got %v", result)
	}
}
