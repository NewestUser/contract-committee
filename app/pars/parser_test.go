package pars

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestVariableCreationAndFunctionPipelining(t *testing.T) {
	useCase := `"customerNumber":"{{$custNum := rndStr | save}}"`;

	funcs := FuncMap{
		"rndStr":rndStr,
		"save": assertSave(t, "$custNum", "gaga"),
	}

	tmpl, _ := NewTemplate("foo").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"gaga"`

	assert.Equal(t, want, got)
}

func TestFunctionExecution(t *testing.T) {
	useCase := `"customerNumber":"{{rndStr}}"`;

	funcs := FuncMap{
		"rndStr":rndStr,
	}

	tmpl, _ := NewTemplate("bar").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"123"`

	assert.Equal(t, want, got)
}


func rndStr() string {
	fmt.Println("================ RND STRING CALLED ==================")
	return "123"
}

func assertSave(t *testing.T, wantName, wantVal string) (func() string) {
	return func() string {
		fmt.Println("================ SAVE CALLED ==================")
		//if wantName != gotName {
		//	t.Errorf("wantName: %v gotName: %v", wantName, gotName)
		//}
		//
		//if wantVal != gotVal {
		//	t.Errorf("wantVal: %v gotVal: %v", wantVal, gotVal)
		//}

		return "gaga"
	}
}
