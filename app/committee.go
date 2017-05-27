package app

type NewTest struct {
	Name string
}

type Test struct {
	Name string
	ID   string
}

type Committee interface {
	RegisterTest(t *NewTest) *Test
}

type persistentCommittee struct {
	repo TestStorage
}

func (c *persistentCommittee) RegisterTest(t *NewTest) *Test {
	id := c.repo.Register(t)
	
	return &Test{Name:t.Name, ID:id}
}

func NewCommittee(repo TestStorage) Committee {
	return &persistentCommittee{repo:repo}
}
