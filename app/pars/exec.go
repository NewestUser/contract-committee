package pars

import (
	"reflect"
	"text/template/parse"
	"fmt"
	"io"
)

type executor struct {
	parseFuncs map[string]interface{}
	varFuncs map[string]interface{}
	w          io.Writer
}



func (e *executor) execFunc(funcName string, varNames []*parse.VariableNode, args []parse.Node, prevCmdResult reflect.Value) reflect.Value {
	//var funcToEval reflect.Value

	parseFunc, parseOk := e.parseFuncs[funcName]
	if (parseOk) {
		funcToEval := reflect.ValueOf(parseFunc)
		return e.evalCall(funcName, funcToEval, args, prevCmdResult)
	}

	varFunc, varOk := e.varFuncs[funcName]
	if (varOk) {
		funcToEval := reflect.ValueOf(varFunc)
		return e.evalVarCall(funcName, funcToEval, varNames[0], args, prevCmdResult)
	}

	panic(fmt.Sprintf("can't find function '%s' to execute", funcName))
}

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

func (e *executor) evalVarCall(funName string, fun reflect.Value, varName *parse.VariableNode, funArgs[]parse.Node, prevFunResult reflect.Value) reflect.Value {
	funType := fun.Type()

	if funType.IsVariadic() {
		panic("variadic functions are not supported")
	}

	numIn := len(funArgs) + 1 // +1 because function name will be passed also as an argument
	numFixed := len(funArgs)
	if (prevFunResult.IsValid()) {
		numIn++
	}

	fmt.Println("numIn", numIn, "numFixed", numFixed)
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
		argv[i] = evalArg(funType.In(i), funArgs[i])
	}

	// Add final value if necessary.
	if prevFunResult.IsValid() {
		t := funType.In(funType.NumIn() - 2)
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
		i++
	}

	varNameArg := reflect.ValueOf(varName.String())
	t := funType.In(funType.NumIn() - 1)
	argv[i] = validateType(varNameArg, t)

	result := fun.Call(argv)
	return result[0]
}


// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}


// validateType guarantees that the value is valid and assignable to the type.
func validateType(value reflect.Value, typ reflect.Type) reflect.Value {
	if !value.IsValid() {
		if typ == nil || canBeNil(typ) {
			// An untyped nil interface{}. Accept as a proper nil value.
			return reflect.Zero(typ)
		}
		panic("invalid value; expected " + typ.String())
	}
	if typ != nil && !value.Type().AssignableTo(typ) {
		if value.Kind() == reflect.Interface && !value.IsNil() {
			value = value.Elem()
			if value.Type().AssignableTo(typ) {
				return value
			}
			// fallthrough
		}
		// Does one dereference or indirection work? We could do more, as we
		// do with method receivers, but that gets messy and method receivers
		// are much more constrained, so it makes more sense there than here.
		// Besides, one is almost always all you need.
		switch {
		case value.Kind() == reflect.Ptr && value.Type().Elem().AssignableTo(typ):
			value = value.Elem()
			if !value.IsValid() {
				panic("dereference of nil pointer of type " + typ.String())
			}
		case reflect.PtrTo(value.Type()).AssignableTo(typ) && value.CanAddr():
			value = value.Addr()
		default:
			panic("wrong type for value; expected but got : " + typ.String() + " " + value.Type().String())
		}
	}
	return value
}

func evalBool(typ reflect.Type, n parse.Node) reflect.Value {
	if n, ok := n.(*parse.BoolNode); ok {
		value := reflect.New(typ).Elem()
		value.SetBool(n.True)
		return value
	}
	panic("expected bool; found " + n.String())
}

func evalInteger(typ reflect.Type, n parse.Node) reflect.Value {
	if n, ok := n.(*parse.NumberNode); ok && n.IsInt {
		value := reflect.New(typ).Elem()
		value.SetInt(n.Int64)
		return value
	}
	panic("expected integer; found " + n.String())
}

func evalUnsignedInteger(typ reflect.Type, n parse.Node) reflect.Value {

	if n, ok := n.(*parse.NumberNode); ok && n.IsUint {
		value := reflect.New(typ).Elem()
		value.SetUint(n.Uint64)
		return value
	}
	panic("expected unsigned integer; found " + n.String())
}

func evalFloat(typ reflect.Type, n parse.Node) reflect.Value {

	if n, ok := n.(*parse.NumberNode); ok && n.IsFloat {
		value := reflect.New(typ).Elem()
		value.SetFloat(n.Float64)
		return value
	}
	panic("expected float; found " + n.String())
}

func evalComplex(typ reflect.Type, n parse.Node) reflect.Value {
	if n, ok := n.(*parse.NumberNode); ok && n.IsComplex {
		value := reflect.New(typ).Elem()
		value.SetComplex(n.Complex128)
		return value
	}
	panic("expected complex; found " + n.String())
}
