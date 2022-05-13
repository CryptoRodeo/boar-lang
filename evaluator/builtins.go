package evaluator

import (
	"fmt"
	"monkey/object"
)

type ArrayErrorFormatter struct {
	FuncName          string
	ArgumentsExpected int
	Arguments         []object.Object
}

var builtins = map[string]*object.Builtin{
	//len()
	"len":   {Fn: __len__},
	"first": {Fn: __first__},
	"last":  {Fn: __last__},
	"rest":  {Fn: __rest__},
	"push":  {Fn: __push__},
	"puts":  {Fn: __puts__},
}

func checkForArrayErrors(formatter ArrayErrorFormatter) object.Object {
	args := formatter.Arguments
	functionName := formatter.FuncName
	argumentsExpected := formatter.ArgumentsExpected
	if len(args) != argumentsExpected {
		return newError("wrong number of arguments, got %d wanted %d", len(args), argumentsExpected)
	}

	if !isArray(args[0]) {
		return newError("argument to `%s` must be ARRAY, got %s", functionName, args[0].Type())
	}

	return NULL
}

func __len__(args ...object.Object) object.Object {
	// len() should only be passed 1 argument
	if len(args) != 1 {
		return newError("wrong number of arguments. got %d, wanted 1", len(args))
	}

	switch arg := args[0].(type) {

	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}

	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}

	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func __first__(args ...object.Object) object.Object {
	err := checkForArrayErrors(ArrayErrorFormatter{FuncName: "first", ArgumentsExpected: 1, Arguments: args})

	if err != NULL {
		return err
	}

	arr := args[0].(*object.Array)

	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func __last__(args ...object.Object) object.Object {
	err := checkForArrayErrors(ArrayErrorFormatter{FuncName: "last", ArgumentsExpected: 1, Arguments: args})

	if err != NULL {
		return err
	}

	arr := args[0].(*object.Array)

	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL

}

/**
- returns a new array containing all the elements of the array passed as argument *except* the first one.
- Similar to the cdr function in Scheme (also similar to tail)
**/
func __rest__(args ...object.Object) object.Object {
	err := checkForArrayErrors(ArrayErrorFormatter{FuncName: "rest", ArgumentsExpected: 1, Arguments: args})

	if err != NULL {
		return err
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])

		return &object.Array{Elements: newElements}
	}

	return NULL

}

/**
- Returns a new array with the pushed element at the end.
- Arrays are immutable in monke-lang, so it doesn't modify the given array
**/
func __push__(args ...object.Object) object.Object {
	err := checkForArrayErrors(ArrayErrorFormatter{FuncName: "push", ArgumentsExpected: 2, Arguments: args})

	if err != NULL {
		return err
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func __puts__(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}
