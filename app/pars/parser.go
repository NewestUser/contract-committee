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

	Parse(tmpl string) (Template, error)

	Execute() string
}

type FuncMap map[string]interface{}

type template struct {
	name       string
	parseFuncs map[string]interface{}

	parsedTree *parse.Tree
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

	t.parsedTree = tree

	return t, nil
}

func (t *template) Execute() string {
	buffer := new(bytes.Buffer)

	root := t.parsedTree.Root

	exec := executor{parseFuncs:t.parseFuncs, w:buffer}

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
		fmt.Println("textNode --")
		exec.w.Write(n.Text)
		return

	case *parse.ActionNode:
		fmt.Println("actionNode --")

		var cmdResult reflect.Value
		for _, cmd := range n.Pipe.Cmds {

			funcName := cmd.Args[0].String()
			funcArgs := cmd.Args[1:]

			//funcToEval := reflect.ValueOf(exec.parseFuncs[funcName])

			cmdResult = exec.execFunc(funcName, funcArgs, cmdResult)
			fmt.Println("funcName", funcName, "funcArgs", funcArgs, "result", cmdResult)

		}

		exec.printValue(cmdResult)

		for i, cmd := range n.Pipe.Decl {
			fmt.Println("pipe decl indent", i, cmd.Ident)
		}
		return
	}
}

func (e *executor) execFunc(funcName string, args []parse.Node, prevCmdResult reflect.Value) reflect.Value {
	funcToEval := reflect.ValueOf(e.parseFuncs[funcName])

	return e.evalCall(funcName, funcToEval, args, prevCmdResult)
}

//
//func (e *executor) execFunc(funcName string, args []parse.Node) reflect.Value {
//	funcToEval := reflect.ValueOf(e.parseFuncs[funcName])
//
//
//	for _, cmd := range pipe.Cmds {
//		value = s.evalCommand(dot, cmd, value) // previous value is this one's final arg.
//		// If the object has type interface{}, dig down one level to the thing inside.
//		if value.Kind() == reflect.Interface && value.Type().NumMethod() == 0 {
//			value = reflect.ValueOf(value.Interface()) // lovely!
//		}
//	}
//	for _, variable := range pipe.Decl {
//		s.push(variable.Ident[0], value)
//	}
//
//
//	return e.evalCall(funcToEval, args)
//}

func (e *executor) printValue(v reflect.Value) {
	pval, ok := printableValue(v)
	if !ok {
		panic("cant print value")
	}

	_, err := fmt.Fprint(e.w, pval)
	if err != nil {
		panic("cant write value")
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

func (e *executor) evalCall(funName string, fun reflect.Value, funArgs[]parse.Node, prevFunResult reflect.Value) reflect.Value {
	funType := fun.Type()

	if funType.IsVariadic() {
		panic("variadic functions are not supported")
	}

	numIn := len(funArgs)
	numFixed := len(funArgs)
	if (prevFunResult.IsValid()) {
		numIn++
	}

	fmt.Println("args", funArgs, "prevResult", prevFunResult)

	if numIn < funType.NumIn() - 1 || !funType.IsVariadic() && numIn != funType.NumIn() {
		panic(fmt.Sprintf("wrong number of args for %s: want %d got %d", funName, funType.NumIn(), numFixed))
	}

	if !goodFunc(funType) {
		panic(fmt.Sprintf("can't call method/function %q with %d results", funName, funType.NumOut()))
	}

	// Build the arg list.
	argv := make([]reflect.Value, numIn)
	// Args must be evaluated.

	i := 0
	for ; i < numIn && i < numFixed; i++ {
		fmt.Println("iterating ================@KEL@J!#KLJ!@KL#@J!LK#J@L!KJ@#KL!")
		argv[i] = evalArg(funType.In(i), funArgs[i])
	}

	// Add final value if necessary.
	if prevFunResult.IsValid() {
		t := funType.In(funType.NumIn() - 1)
		//if funType.IsVariadic() {
		//	panic()
		//if numIn-1 < numFixed {
		//	 The added final argument corresponds to a fixed parameter of the function.
		//	 Validate against the type of the actual parameter.
		//t = typ.In(numIn - 1)
		//} else {
		//	 The added final argument corresponds to the variadic part.
		//	 Validate against the type of the elements of the variadic slice.
		//t = t.Elem()
		//}
		//}
		argv[i] = validateType(prevFunResult, t)
	}

	result := fun.Call(argv)
	return result[0]
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

// evalCall executes a function or method call. If it's a method, fun already has the receiver bound, so
// it looks just like a function call. The arg list, if non-nil, includes (in the manner of the shell), arg[0]
// as the function itself.
//func evalCall(dot, fun reflect.Value, node parse.Node, name string, args []parse.Node, final reflect.Value) reflect.Value {
//	if args != nil {
//		args = args[1:] // Zeroth arg is function name/node; not passed to function.
//	}
//	typ := fun.Type()
//	numIn := len(args)
//	if final.IsValid() {
//		numIn++
//	}
//	numFixed := len(args)
//	if typ.IsVariadic() {
//		numFixed = typ.NumIn() - 1 // last arg is the variadic one.
//		if numIn < numFixed {
//			s.errorf("wrong number of args for %s: want at least %d got %d", name, typ.NumIn()-1, len(args))
//		}
//	} else if numIn < typ.NumIn()-1 || !typ.IsVariadic() && numIn != typ.NumIn() {
//		s.errorf("wrong number of args for %s: want %d got %d", name, typ.NumIn(), len(args))
//	}
//	if !goodFunc(typ) {
//		// TODO: This could still be a confusing error; maybe goodFunc should provide info.
//		s.errorf("can't call method/function %q with %d results", name, typ.NumOut())
//	}
//	// Build the arg list.
//	argv := make([]reflect.Value, numIn)
//	// Args must be evaluated. Fixed args first.
//	i := 0
//	for ; i < numFixed && i < len(args); i++ {
//		argv[i] = s.evalArg(dot, typ.In(i), args[i])
//	}
//	// Now the ... args.
//	if typ.IsVariadic() {
//		argType := typ.In(typ.NumIn() - 1).Elem() // Argument is a slice.
//		for ; i < len(args); i++ {
//			argv[i] = s.evalArg(dot, argType, args[i])
//		}
//	}
//	// Add final value if necessary.
//	if final.IsValid() {
//		t := typ.In(typ.NumIn() - 1)
//		if typ.IsVariadic() {
//			if numIn-1 < numFixed {
//				// The added final argument corresponds to a fixed parameter of the function.
//				// Validate against the type of the actual parameter.
//				t = typ.In(numIn - 1)
//			} else {
//				// The added final argument corresponds to the variadic part.
//				// Validate against the type of the elements of the variadic slice.
//				t = t.Elem()
//			}
//		}
//		argv[i] = s.validateType(final, t)
//	}
//	result := fun.Call(argv)
//	// If we have an error that is not nil, stop execution and return that error to the caller.
//	if len(result) == 2 && !result[1].IsNil() {
//		s.at(node)
//		s.errorf("error calling %s: %s", name, result[1].Interface().(error))
//	}
//	return result[0]
//}


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
	//case reflect.Bool:
	//	return s.evalBool(typ, n)
	//case reflect.Complex64, reflect.Complex128:
	//	return s.evalComplex(typ, n)
	//case reflect.Float32, reflect.Float64:
	//	return s.evalFloat(typ, n)
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//	return s.evalInteger(typ, n)
	//case reflect.Interface:
	//	if typ.NumMethod() == 0 {
	//		return s.evalEmptyInterface(dot, n)
	//	}
	case reflect.String:
		return evalString(typ, n)
	//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	//	return s.evalUnsignedInteger(typ, n)
	}
	panic(fmt.Sprintf("can't handle %s for arg of type %s", n, typ))
}



//func evalBool(typ reflect.Type, n parse.Node) reflect.Value {
//	s.at(n)
//	if n, ok := n.(*parse.BoolNode); ok {
//		value := reflect.New(typ).Elem()
//		value.SetBool(n.True)
//		return value
//	}
//	s.errorf("expected bool; found %s", n)
//	panic("not reached")
//}

func evalString(typ reflect.Type, n parse.Node) reflect.Value {
	//s.at(n)
	if n, ok := n.(*parse.StringNode); ok {
		value := reflect.New(typ).Elem()
		value.SetString(n.Text)
		return value
	}

	fmt.Println(" EVAL STRING WAS CALLED YAYYY ==============")
	panic(fmt.Sprintf("expected string; found %s", n))
}
//
//func evalInteger(typ reflect.Type, n parse.Node) reflect.Value {
//	s.at(n)
//	if n, ok := n.(*parse.NumberNode); ok && n.IsInt {
//		value := reflect.New(typ).Elem()
//		value.SetInt(n.Int64)
//		return value
//	}
//	panic(fmt.Sprintf("expected integer; found %s", n))
//}
//
//func evalUnsignedInteger(typ reflect.Type, n parse.Node) reflect.Value {
//	s.at(n)
//	if n, ok := n.(*parse.NumberNode); ok && n.IsUint {
//		value := reflect.New(typ).Elem()
//		value.SetUint(n.Uint64)
//		return value
//	}
//	panic(fmt.Sprintf("expected unsigned integer; found %s", n))
//}
//
//func evalFloat(typ reflect.Type, n parse.Node) reflect.Value {
//	s.at(n)
//	if n, ok := n.(*parse.NumberNode); ok && n.IsFloat {
//		value := reflect.New(typ).Elem()
//		value.SetFloat(n.Float64)
//		return value
//	}
//	panic(fmt.Sprintf("expected float; found %s", n))
//}

func printNodeStuff(node parse.Node) {

	switch n := node.(type) {
	case *parse.TextNode:
		fmt.Println("textNode--")
		fmt.Println(string(n.Text))
		return

	case *parse.ActionNode:
		fmt.Println("actionNode --")
		for i, cmd := range n.Pipe.Cmds {
			fmt.Println("pipe cmd", i, cmd)
			fmt.Println("pipe cmd args", i, cmd.Args)
		}

		for i, cmd := range n.Pipe.Decl {
			fmt.Println("pipe decl indent", i, cmd.Ident)
		}
		return
	}

	fmt.Println("NODE:", node)
	fmt.Println("NODE Type:", node.Type())
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