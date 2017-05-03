package pars

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"strings"
)

func TestVariableCreationAndFunctionPipelining(t *testing.T) {
	useCase := `"customerNumber":"{{$custNum := rndStr | toUpper}}"`;

	funcs := FuncMap{
		"rndStr":fooString,
		"toUpper": toUpperCase,
	}

	tmpl, _ := NewTemplate("any").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"FOO"`

	assert.Equal(t, want, got)
}

//func TestVariableCreation(t *testing.T) {
//	useCase := `{{$myVar := rndStr}}`
//
//	execFuncs := FuncMap{
//		"rndStr":rndStr,
//	}
//	
//
//}

func TestFunctionExecution(t *testing.T) {
	useCase := `"customerNumber":"{{rndStr}}"`;

	funcs := FuncMap{
		"rndStr":fooString,
	}

	tmpl, _ := NewTemplate("bar").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"foo"`

	assert.Equal(t, want, got)
}

func fooString() string {
	fmt.Println("================ RND STRING CALLED ==================")
	return "foo"
}

func toUpperCase(str string) string {
	return strings.ToUpper(str)
}

func assertSave(t *testing.T, wantName, wantVal string) (func() string) {
	return func() string {
		fmt.Println("================ SAVE CALLED ==================")
		//if wantName != wantName {
		//	t.Errorf("wantName: %v gotName: %v", wantName, gotName)
		//}
		//
		//if wantVal != gotVal {
		//	t.Errorf("wantVal: %v gotVal: %v", wantVal, gotVal)
		//}

		return "gaga"
	}
}
