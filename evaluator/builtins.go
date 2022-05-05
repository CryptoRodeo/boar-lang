package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	//len()
	"len":   {Fn: __len__},
	"first": {Fn: __first__},
	"last":  {Fn: __last__},
}

func checkForArrayErrors(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got %d wanted 1", len(args))
	}

	if !isArray(args[0]) {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
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
	checkForArrayErrors(args...)

	arr := args[0].(*object.Array)

	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func __last__(args ...object.Object) object.Object {
	checkForArrayErrors(args...)

	arr := args[0].(*object.Array)

	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL

}
