package assert

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

type Runner interface {
	
}

func NewRunner()  {
	
}
