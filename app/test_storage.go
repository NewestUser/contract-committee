package app

type SuiteStorage interface {
	Register(t *NewSuite) string

	Exists(name string) bool

	Find(id string) *Suite
}

type CaseStorage interface {
	Register(suiteID string, c *NewCase) string
}