package pars

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"fmt"
)

func TestVariableCreationAndFunctionPipelining(t *testing.T) {
	useCase := `"customerNumber":"{{$custNum := rndStr | toUpper}}"`;

	funcs := FuncMap{
		"rndStr":returnStr("bar"),
		"toUpper": toUpperCase,
	}

	tmpl, _ := NewTemplate("any").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"BAR"`

	assert.Equal(t, want, got)
}

func TestFunctionExecution(t *testing.T) {
	useCase := `"customerNumber":"{{rndStr}}"`;

	funcs := FuncMap{
		"rndStr":returnStr("foo"),
	}

	tmpl, _ := NewTemplate("bar").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"customerNumber":"foo"`

	assert.Equal(t, want, got)
}

func TestSimpleFunctionExecutionWithInt(t *testing.T) {
	useCase := `"{{rndInt | double}}"`;

	funcs := FuncMap{
		"rndInt":returnInt(112),
		"double":doubleInt,
	}

	tmpl, _ := NewTemplate("bar").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"224"`

	assert.Equal(t, want, got)
}

func TestFunctionExecutionWithStringArguments(t *testing.T) {
	useCase := `"banana-{{concat "lo" "ve"}}"`

	funcs := FuncMap{
		"concat":concatStrings,
	}

	tmpl, e := NewTemplate("aaa").Funcs(funcs).Parse(useCase)

	fmt.Println(e)
	got := tmpl.Execute()
	want := `"banana-love"`

	assert.Equal(t, want, got)
}

func TestFunctionExecutionWithIntArguments(t *testing.T) {
	useCase := `"banana-{{add 1 3}}"`

	funcs := FuncMap{
		"add":addInts,
	}

	tmpl, _ := NewTemplate("aaa").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := `"banana-4"`

	assert.Equal(t, want, got)
}

func TestFunctionExecutionWithMixedArguments(t *testing.T) {
	useCase := `{{mixConcat "aaa" 73 | toUpper}}`

	funcs := FuncMap{
		"mixConcat" : strAndIntConcat,
		"toUpper": toUpperCase,
	}

	tmpl, _ := NewTemplate("does not matter").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := "AAA73"

	assert.Equal(t, want, got)
}

func TestFunctionExecutionWithBoolArgument(t *testing.T) {
	useCase := `{{invert true}}`

	funcs := FuncMap{
		"invert" : invert,
	}

	tmpl, _ := NewTemplate("does not matter").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := "false"

	assert.Equal(t, want, got)
}

func TestFunctionExecutionWithFloatArguments(t *testing.T) {
	useCase := `{{add 2.4 3.4}}`

	funcs := FuncMap{
		"add": addFloats,
	}

	tmpl, _ := NewTemplate("doesn't matter").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := "5.8"

	assert.Equal(t, want, got)
}

func TestVariableCreation(t *testing.T) {
	useCase := `{{$foo := rndStr | save}}`;

	funcs := FuncMap{
		"rndStr":returnStr("bar"),
	}

	var varName string
	var value string

	varFuncs := VarFuncMap{
		"save":func(val string, name string) string {
			value = val
			varName = name
			return value
		},
	}

	tmpl, _ := NewTemplate("tmpl").Funcs(funcs).VarFuncs(varFuncs).Parse(useCase)
	got := tmpl.Execute()

	wantResult := "bar"
	wantVal := "bar"
	wantName := "$foo"

	assert.Equal(t, wantResult, got)
	assert.Equal(t, wantName, varName)
	assert.Equal(t, wantVal, value)
}

func TestPipeliningConcreteValue(t *testing.T) {
	useCase := `{{true | invert}}`

	funcs := FuncMap{
		"invert" : invert,
	}

	tmpl, _ := NewTemplate("does not matter").Funcs(funcs).Parse(useCase)

	got := tmpl.Execute()
	want := "false"

	assert.Equal(t, want, got)
}

func returnStr(want string) (func() string) {
	return func() string {
		return want
	}
}

func toUpperCase(str string) string {
	return strings.ToUpper(str)
}

func returnInt(want int) (func() int) {
	return func() int {
		return want
	}
}

func doubleInt(val int) int {
	return val * 2
}

func concatStrings(argOne, argTwo string) string {
	return argOne + argTwo
}

func addInts(i1, i2 int) int {
	return i1 + i2
}

func strAndIntConcat(str string, i int) string {
	return fmt.Sprintf("%v%v", str, i)
}

func invert(b bool) bool {
	return !b
}

func addFloats(f1, f2 float64) float64 {
	return f1 + f2
}