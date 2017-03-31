package pars

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestVariableCreationAndFunctionPipelining(t *testing.T) {
	useCase := `"customerNumber":"{{$custNum := rndStr | save}}"`;

	funcs := FuncMap{
		"rndStr":rndStr,
		"save": assertSave(t, "$custNum", "123"),
	}

	tmpl, _ := NewTemplate("foo").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"123"`

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

func TestNewTemplate(t *testing.T) {



}

func rndStr() string {
	return "123"
}

func assertSave(t *testing.T, wantName, wantVal string) (func(string, string) string) {
	return func(gotName, gotVal string) string {
		if wantName != gotName {
			t.Errorf("wantName: %v gotName: %v", wantName, gotName)
		}

		if wantVal != gotVal {
			t.Errorf("wantVal: %v gotVal: %v", wantVal, gotVal)
		}

		return gotVal
	}
}
