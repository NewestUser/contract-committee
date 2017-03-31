package datastore

type Collection interface {
	// Insert inserts a new document in the datastore.
	// In case or error it returns the original error that was occurred
	Insert(doc interface{}) error

	// Upsert is creating or updating existing document. In case when document exists, then
	// update operation is performed, otherwise update operation is used as insert.
	Upsert(query interface{}, update interface{}) error

	// Aggregate executes the provided pipeline into provided result.
	Aggregate(pipeline interface{}, result interface{}) error

	// FindSorted finds all items from the named collection in sorted order.
	FindSorted(query interface{}, sort string, limit int, result interface{}) error

	// Find finds all items for given query in named collection
	Find(query interface{}, result interface{}) error
}

type DB interface {
	C(name string) Collection
}
