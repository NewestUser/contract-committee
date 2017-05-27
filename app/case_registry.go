package app

type Given struct {
	URL    string
	Method string
	Body   interface{}
}

type Assertion struct {
	StatusCode int
	Body       interface{}
}

type NewCase struct {
	GivenReq   *Given
	AssertResp *Assertion
}

type Case struct {
	ID string
	*NewCase
}

type CaseRegistry interface {
	RegisterCase(suiteID string, c *NewCase) *Case
}

func NewCaseRegistry(s CaseStorage) CaseRegistry {
	return &persistentCaseRegistry{storage:s}
}

type persistentCaseRegistry struct {
	storage CaseStorage
}

func (r *persistentCaseRegistry) RegisterCase(suiteID string, c *NewCase) *Case {

	id := r.storage.Register(suiteID, c)

	return &Case{ID:id, NewCase:c}
}