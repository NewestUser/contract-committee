package assert

import (
	"github.com/newestuser/contract-committee/app/pars"
	"fmt"
)

type Validator interface {
	Valid(c *NewCase) error
}

type caseValidator struct {

}

func (v *caseValidator) Valid(c *NewCase) error {
	_, err := pars.NewTemplate("foo").Parse(fmt.Sprintf("%s", c.GivenReq.Body))

	return err
}