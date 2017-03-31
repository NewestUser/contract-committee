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

	treeMap, err := parse.Parse(t.name, tmpl, leftDelim, rightDelim, t.parseFuncs)

	if err != nil {
		return nil, err
	}

	tree := treeMap[t.name]
	root := tree.Root

	//for k, v := range treeMap {
	//	fmt.Println("key:", k)
	//	fmt.Println("val:", v)
		fmt.Println("============================================")
		//root := v.Root
		fmt.Println("root:", root)
		fmt.Println("============================================")

		for i, n := range root.Nodes {
			fmt.Println("n", i, ":", n)
			//fmt.Println("type", n.Type())
			printNodeStuff(n)
		}
	//}


	return t, nil
}

func (t *template) Execute() string {
	buffer := new(bytes.Buffer)



	return buffer.String()
}

func (t *template) Name() string  {
	return t.name
}



func printNodeStuff(node parse.Node) {
	switch n := node.(type) {
	case *parse.ActionNode:
		//fmt.Println("line", n.Line)
		//fmt.Println("pos", n.Pos)
		//fmt.Println("pipe", n.Pipe)

		for i, cmd := range n.Pipe.Cmds{
			fmt.Println("pipe cmd",i,cmd)
			fmt.Println("pipe cmd args",i,cmd.Args)
		}

		for i, cmd := range n.Pipe.Decl{
			fmt.Println("pipe decl indent",i,cmd.Ident)
		}
		return
	}
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