package pars

import (
	"reflect"
	"fmt"
	"text/template/parse"
	"bytes"
)

const (
	leftDelim = "{{"
	rightDelim = "}}"
)


type Template interface {
	Name() string

	Funcs(funcMap FuncMap) Template

	Parse(tmpl string) (Template, error)

	Execute() string

	ErrorContext(node parse.Node) (location, context string)
}

type FuncMap map[string]interface{}

type template struct {
	name           string
	parseFuncs     map[string]interface{}

	parsedTemplate *parse.Tree
}

func NewTemplate(name string) Template {
	return &template{name:name, parseFuncs:make(map[string]interface{})}
}

func (t *template) Funcs(funcs FuncMap) Template {
	t.parseFuncs = funcs
	return t
}

func (t *template) Parse(tmpl string) (Template, error) {

	tree, err := parse.Parse(t.name, tmpl, leftDelim, rightDelim, t.parseFuncs)

	if err != nil {
		return nil, err
	}

	t.parsedTemplate = tree[t.name]
	//
	//var buffer bytes.Buffer
	//
	//startIdnex := strings.Index(tmpl, leftDelim)
	//endIdnex := strings.Index(tmpl, rightDelim)
	//
	//buffer.WriteString(tmpl[0:startIdnex])
	//
	//equation := tmpl[startIdnex:endIdnex]
	//noDelimEquation := tmpl[startIdnex + len(leftDelim) : endIdnex + len(rightDelim)]
	//
	//t.process(noDelimEquation)
	//
	//buffer.WriteString(equation)
	//buffer.WriteString(tmpl[endIdnex:])
	//
	//t.parsedTemplate = buffer

	return t, nil
}

func (t *template) Execute() string {
	buffer := new(bytes.Buffer)

	type foo struct {

	}

	value := reflect.ValueOf(&foo{})
	state := &state{
		tmpl: t,
		wr:   buffer,
		vars: []variable{{"$", value}},
	}
	if t.parsedTemplate == nil || t.parsedTemplate.Root == nil {
		state.errorf("%q is an incomplete or empty template%s", t.Name(), "dsadsadsa")
	}
	state.walk(value, t.parsedTemplate.Root)


	return buffer.String()
}

func (t *template) Name() string  {
	return t.name
}

func (t *template) ErrorContext(node parse.Node) (location, context string)  {

	//pos := int(n.Position())
	//tree := n.tree()
	//if tree == nil {
	//	tree = t
	//}
	//text := tree.text[:pos]
	//byteNum := strings.LastIndex(text, "\n")
	//if byteNum == -1 {
	//	byteNum = pos // On first line.
	//} else {
	//	byteNum++ // After the newline.
	//	byteNum = pos - byteNum
	//}
	//lineNum := 1 + strings.Count(text, "\n")
	//context = n.String()
	//if len(context) > 20 {
	//	context = fmt.Sprintf("%.20s...", context)
	//}
	return fmt.Sprintf("foooo"), context
	//return
}

func (t *template) process(equation string) {
	//noDelimEqu := stripDelims(equation)
	// t.parseFuncs[equation]

}

// length returns the length of the item, with an error if it has no defined length.
func length(item interface{}) (int, error) {
	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return 0, fmt.Errorf("len of untyped nil")
	}
	v, isNil := indirect(v)
	if isNil {
		return 0, fmt.Errorf("len of nil pointer")
	}
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len(), nil
	}
	return 0, fmt.Errorf("len of type %s", v.Type())
}


// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}