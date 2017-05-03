package pars

import (
	"reflect"
	"text/template/parse"
	"fmt"
	"io"
)

type executor struct {
	parseFuncs map[string]interface{}
	w          io.Writer
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

//func evalString(typ reflect.Type, n parse.Node) reflect.Value {
//	if n, ok := n.(*parse.StringNode); ok {
//		value := reflect.New(typ).Elem()
//		value.SetString(n.Text)
//		return value
//	}
//	panic("expected string; found " + n.String())
//}

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

func evalEmptyInterface(dot reflect.Value, n parse.Node) reflect.Value {

	switch n := n.(type) {
	case *parse.BoolNode:
		return reflect.ValueOf(n.True)
	case *parse.DotNode:
		return dot
	case *parse.FieldNode:
		fmt.Println("implement FieldNode")
	//return evalFieldNode(dot, n, nil, zero)
	case *parse.IdentifierNode:
		fmt.Println("implement IdentifierNode")
	//return evalFunction(dot, n, n, nil, zero)
	case *parse.NilNode:
		// NilNode is handled in evalArg, the only place that calls here.
		panic("evalEmptyInterface: nil (can't happen)")
	case *parse.NumberNode:
		fmt.Println("implemenet NumberNode")
	//return idealConstant(n)
	case *parse.StringNode:
		return reflect.ValueOf(n.Text)
	case *parse.VariableNode:
		fmt.Println("implement VariableNode")
	//return evalVariableNode(dot, n, nil, zero)
	case *parse.PipeNode:
		fmt.Println("implement PipeNode")
	//return evalPipeline(dot, n)
	}
	panic("can't handle assignment of " + n.String() + " interface argument")
}
