package evaluator

import (
	"fmt"
	"monkey/object"
)

type ErrorFormatter struct {
	FuncName          string
	ArgumentsExpected int //minimum arguments expected
	Arguments         []object.Object
}

var builtins = map[string]*object.Builtin{
	//len()
	"len":      {Fn: __len__},
	"first":    {Fn: __first__},
	"last":     {Fn: __last__},
	"rest":     {Fn: __rest__},
	"push":     {Fn: __push__},
	"puts":     {Fn: __puts__},
	"delete":   {Fn: __delete__},
	"valuesAt": {Fn: __valuesAt__},
	"toArray":  {Fn: __toArray__},
	"dig":      {Fn: __dig__},
}

func checkForArrayErrors(formatter ErrorFormatter) object.Object {
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
	err := checkForArrayErrors(ErrorFormatter{FuncName: "first", ArgumentsExpected: 1, Arguments: args})

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
	err := checkForArrayErrors(ErrorFormatter{FuncName: "last", ArgumentsExpected: 1, Arguments: args})

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
	err := checkForArrayErrors(ErrorFormatter{FuncName: "rest", ArgumentsExpected: 1, Arguments: args})

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
	err := checkForArrayErrors(ErrorFormatter{FuncName: "push", ArgumentsExpected: 2, Arguments: args})

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

func __delete__(args ...object.Object) object.Object {
	err := checkForHashErrors(ErrorFormatter{FuncName: "delete", ArgumentsExpected: 2, Arguments: args})

	if err != NULL {
		return err
	}

	// First argument must be a hash
	hash := args[0].(*object.Hash)

	// The remaining arguments should be valid hash keys.
	// Loop through them and null their values
	for _, arg := range args[1:] {
		hashKey, ok := arg.(object.Hashable)

		if !ok {
			return newError("Unusable value as hash key: %s", arg.Type())
		}

		hash.Pairs[hashKey.HashKey()] = object.HashPair{Key: NULL, Value: NULL}
	}

	return hash
}

func __valuesAt__(args ...object.Object) object.Object {
	err := checkForHashErrors(ErrorFormatter{FuncName: "valuesAt", ArgumentsExpected: 2, Arguments: args})

	if err != NULL {
		return err
	}

	// First argument must be a hash
	hash := args[0].(*object.Hash)

	// Create array object to store object values at x key
	arr := &object.Array{}

	// The remaining arguments should be valid hash keys.
	for _, arg := range args[1:] {
		hashKey, ok := arg.(object.Hashable)

		if !ok {
			return newError("Unusable value as hash key: %s", arg.Type())
		}

		// Grab the value at said key, append to array
		arr.Elements = append(arr.Elements, hash.Pairs[hashKey.HashKey()].Value)
	}

	return arr
}

func __toArray__(args ...object.Object) object.Object {
	err := checkForHashErrors(ErrorFormatter{FuncName: "valuesAt", ArgumentsExpected: 1, Arguments: args})

	if err != NULL {
		return err
	}

	// First argument must be a hash
	hash := args[0].(*object.Hash)

	// Create array object to store object values at x key
	arr := &object.Array{}

	for _, pair := range hash.Pairs {
		arr.Elements = append(arr.Elements, pair.Key, pair.Value)
	}

	return arr
}

func __dig__(args ...object.Object) object.Object {
	err := checkForHashErrors(ErrorFormatter{FuncName: "valuesAt", ArgumentsExpected: 2, Arguments: args})

	if err != NULL {
		return err
	}

	// First argument must be a hash
	hash := args[0].(*object.Hash)

	if len(args) == 1 {
		return hash
	}

	hashKey, ok := args[1].(object.Hashable)

	if !ok {
		return newError("Unusable value as hash key: %s", args[1].Type())
	}

	extracted, exists := hash.Pairs[hashKey.HashKey()]

	if exists {
		// if we only have 2 args (someInnerHash, key), and we've found the value exists then return it
		if len(args) == 2 {
			return extracted.Value
		}

		newArgs := []object.Object{}
		newArgs = append(newArgs, extracted.Value)
		newArgs = append(newArgs, args[2:]...)

		return __dig__(newArgs...)
	}

	return nil
}

func checkForHashErrors(formatter ErrorFormatter) object.Object {
	args, functionName, argumentsExpected := formatter.Arguments, formatter.FuncName, formatter.ArgumentsExpected

	if len(args) < argumentsExpected {
		return newError("wrong number of arguments, got %d wanted %d", len(args), argumentsExpected)
	}

	if !isHash(args[0]) {
		return newError("argument to `%s` must be HASH, got %s", functionName, args[0].Type())
	}

	return NULL
}
