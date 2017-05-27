package app

type NewSuite struct {
	Name string
}

type Suite struct {
	Name string
	ID   string
}

type SuiteInteractor interface {
	RegisterSuite(s *NewSuite) (*Suite, error)

	FindSuite(id string) (*Suite, error)
}

type persistentSuite struct {
	repo SuiteStorage
}

func (c *persistentSuite) RegisterSuite(s *NewSuite) (*Suite, error) {

	if (c.repo.Exists(s.Name)) {
		return nil, BadReqErr("Test with name %s already exists", s.Name)
	}

	id := c.repo.Register(s)

	return &Suite{Name:s.Name, ID:id}, nil
}

func (c *persistentSuite) FindSuite(id string) (*Suite, error) {

	if t := c.repo.Find(id); t != nil {
		return t, nil
	}

	return nil, NotFoundErr("Test with id %s can not be found", id)
}

func NewSuiteInteractor(repo SuiteStorage) SuiteInteractor {
	return &persistentSuite{repo:repo}
}
