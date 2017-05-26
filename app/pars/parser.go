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

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

type Template interface {
	Name() string

	Funcs(funcMap FuncMap) Template

	VarFuncs(varFuncMap VarFuncMap) Template

	Parse(tmpl string) (Template, error)

	Execute() string
}

type FuncMap map[string]interface{}

type VarFuncMap map[string]interface{}

type template struct {
	name       string
	parseFuncs FuncMap
	varFuncs   VarFuncMap

	parsedTree *parse.Tree
}

func NewTemplate(name string) Template {
	return &template{name:name, parseFuncs:make(map[string]interface{}), varFuncs:make(map[string]interface{})}
}

func (t *template) Funcs(funcs FuncMap) Template {
	t.parseFuncs = funcs
	return t
}

func (t *template) VarFuncs(varFuncs VarFuncMap) Template {
	t.varFuncs = varFuncs
	return t
}

func (t *template) Parse(tmpl string) (Template, error) {
	treeMap, err := parse.Parse(t.name, tmpl, leftDelim, rightDelim, t.parseFuncs, t.varFuncs)

	if err != nil {
		return nil, err
	}

	tree := treeMap[t.name]

	t.parsedTree = tree

	fmt.Println("parse called")
	return t, nil
}

func (t *template) Execute() string {
	buffer := new(bytes.Buffer)

	root := t.parsedTree.Root

	exec := executor{parseFuncs:t.parseFuncs, varFuncs:t.varFuncs, w:buffer}

	for _, n := range root.Nodes {
		t.writeNode(exec, n)
	}

	return buffer.String()
}

func (t *template) Name() string {
	return t.name
}

func (t *template) writeNode(exec executor, node parse.Node) {

	switch n := node.(type) {
	case *parse.TextNode:
		exec.w.Write(n.Text)
		return

	case *parse.ActionNode:
		var cmdResult reflect.Value

		for _, cmd := range n.Pipe.Cmds {
			funcName := cmd.Args[0].String()
			funcArgs := cmd.Args[1:]

			fmt.Println(n.Pipe.Decl)
			fmt.Println("funcName", funcName, "funcArgs", funcArgs, "cmdResult", cmdResult)
			cmdResult = exec.execFunc(funcName, n.Pipe.Decl, funcArgs, cmdResult)
		}

		exec.printValue(cmdResult)

		for i, cmd := range n.Pipe.Decl {
			fmt.Println("pipe decl indent", i, cmd.Ident)
		}

		return
	}
}


// printableValue returns the, possibly indirected, interface value inside v that
// is best for a call to formatted printer.
func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "<no value>", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PtrTo(v.Type()).Implements(errorType) || reflect.PtrTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}


// goodFunc reports whether the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
	// We allow functions with 1 result.
	switch {
	case typ.NumOut() == 1:
		return true
	}
	return false
}

func evalArg(typ reflect.Type, n parse.Node) reflect.Value {
	switch n.(type) {
	case *parse.NilNode:
		if canBeNil(typ) {
			return reflect.Zero(typ)
		}
		panic(fmt.Sprintf("cannot assign nil to %s", typ))
	//case *parse.FieldNode:
	//	return s.validateType(s.evalFieldNode(dot, arg, []parse.Node{n}, zero), typ)
	//case *parse.VariableNode:
	//	return s.validateType(s.evalVariableNode(dot, arg, nil, zero), typ)
	//case *parse.PipeNode:
	//	return s.validateType(s.evalPipeline(dot, arg), typ)
	//case *parse.IdentifierNode:
	//	return s.validateType(s.evalFunction(dot, arg, arg, nil, zero), typ)
	//case *parse.ChainNode:
	//	return s.validateType(s.evalChainNode(dot, arg, nil, zero), typ)
	}
	switch typ.Kind() {
	case reflect.Bool:
		return evalBool(typ, n)
	case reflect.Complex64, reflect.Complex128:
		return evalComplex(typ, n)
	case reflect.Float32, reflect.Float64:
		return evalFloat(typ, n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return evalInteger(typ, n)
	//case reflect.Interface:
	//	if typ.NumMethod() == 0 {
	//		return s.evalEmptyInterface(dot, n)
	//	}
	case reflect.String:
		return evalString(typ, n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return evalUnsignedInteger(typ, n)
	}
	panic(fmt.Sprintf("can't handle %s for arg of type %s", n, typ))
}

func evalString(typ reflect.Type, n parse.Node) reflect.Value {
	//s.at(n)
	if n, ok := n.(*parse.StringNode); ok {
		value := reflect.New(typ).Elem()
		value.SetString(n.Text)
		return value
	}

	panic(fmt.Sprintf("expected string; found %s", n))
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