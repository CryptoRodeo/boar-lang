package evaluator

import (
	"boar/ast"
	"boar/object"
	"fmt"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	//statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

		// a single statement
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	//expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.LetStatement:
		// evaluate the value
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		// assign the value to the identifier: let x = 0
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	/**
	There is no difference between every new boolean we encounter.
	Instead of creating a new instance every time we encounter true or false lets
	just keep referencing the same ones (TRUE, FALSE)
	**/
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		// the operand
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		// now evaluate the operand with the operator
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)

		if isError(left) {
			return left
		}

		if isError(right) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		// note: the env set here is the env/scope the function was defined in
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		// left -> the expression using the index operator: a[0], arr[3]
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		// The index itself
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.IndexAssignment:
		// left -> The expression using the index operator: hash[a], arr[2+2], etc
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		// Evaluate the index expression: [2+2], [10 < 5], etc
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		// Evaluate the value: 2+2, "hello world", etc
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		return evalIndexAssignment(left, index, value)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.InternalFunctionCall:
		// someArr, someHash
		caller_ident := Eval(node.CallerIdentifier, env)

		if isError(caller_ident) {
			return caller_ident
		}
		// .pop(), .delete(), etc
		func_ident := Eval(node.FunctionIdentifier, env)

		if isError(func_ident) {
			return caller_ident
		}
		// (1,2,3), ("a", "b", "c"), etc
		args := evalExpressions(node.Arguments, env)

		for _, arg := range args {
			if isError(arg) {
				return nil
			}
		}
		//  ( someArr/someHash, (1,2,3) )
		newArgs := append([]object.Object{caller_ident}, args...)

		// call the function as usual builtInFunc(objectIdentifier, args)
		return applyFunction(func_ident, newArgs)

	case *ast.AssignmentExpression:
		// x, y, someIdentifier
		// lets make sure this value exists (either in this scope or any existing outer scopes)

		_, exists := env.Get(node.Name.Value)
		if !exists {
			return newError(`Identifier "%s" not found`, node.Name.Value)
		}

		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.ForLoopStatement:
		// Lets set the counter var in the env
		counterVar := Eval(node.CounterVar, env)

		if isError(counterVar) {
			return counterVar
		}

		// Now lets make sure it actually got set / exists
		counterVar, exists := env.Get(node.CounterVar.Name.Value)

		if !exists {
			return counterVar
		}

		// Lets make sure we have a valid for loop condition
		loopCondition := Eval(node.LoopCondition, env)

		if isError(loopCondition) {
			return loopCondition
		}

		// Now run the for loop.
		return applyForLoop(node, env)
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		// if we encounter a return value, immediately return the unwrapped value
		// note: we don't return an object.ReturnValue when encountering it, only the value its wrapping
		case *object.ReturnValue:
			return result.Value
		// if we encounter an error, return immediately
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	//!true
	case TRUE:
		return FALSE
	//!false
	case FALSE:
		return TRUE
	//!null
	case NULL:
		return TRUE
	//!x, !5, etc.
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	//extract value from *object.Integer via type assertion
	value := right.(*object.Integer).Value
	// return integer object with negated value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case bothAreIntegers(left, right):
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case bothAreStrings(left, right):
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func bothAreIntegers(a, b object.Object) bool {
	return isInteger(a) && isInteger(b)
}

func bothAreStrings(a, b object.Object) bool {
	return isString(a) && isString(b)
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	// extract values
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		// if the result is an *object.ReturnValue, return it without unwrapping its .Value
		// and stop the execution in a potential outer block statement.
		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// check if value exists in env
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("identifier not found: " + node.Value)
}

// evaluate expressions (left to right)
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		// if err, stop evaluation, return error
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {

	// check if its a regular function
	case *object.Function:
		arr, isArray := args[0].(*object.Array)
		// Possibly a #map call
		if isArray && len(args) == 1 {
			return applyMapCall(arr, fn)
		}
		// create the inner function scope
		extendedEnv := extendFunctionEnv(fn, args)
		//evalute the function body with the inner scope
		evaluated := Eval(fn.Body, extendedEnv)
		// if the object has a return value, return that value
		// else, return the object.
		return unwrapReturnValue(evaluated)

	// return the built in function, pass args
	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func applyMapCall(arr *object.Array, fn *object.Function) object.Object {
	res := &object.Array{}
	for _, val := range arr.Elements {
		args := []object.Object{val}
		// create an inner function scope for each function call
		extendedEnv := extendFunctionEnv(fn, args)
		// Evaluate
		evaluated := Eval(fn.Body, extendedEnv)
		// Add result to the array
		res.Elements = append(res.Elements, unwrapReturnValue(evaluated))
	}
	return res
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	// create inner function scope
	env := object.NewEnclosedEnvironment(fn.Env)
	// bind the arguments used in the function call to the parameter names in that inner function scope
	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}

	return env
}

/**
returns the unwrapped *object.ReturnValue if it exists.
This is necessary because:
- otherwise a return statement would bubble up through several functions
  - this would stop the evaluation in all of them
- we only want to stop the evaluation of the last called function's body

This is why we need to unwrap it, so evalBlockStatement wont stop evaluating statements in outer functions.
**/
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	// currently we only support concatenation: (string + string)
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}

}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case isArray(left) && isInteger(index):
		return evalArrayIndexExpression(left, index)
	case isHash(left):
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	// total number of elements in the current array
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalIndexAssignment(indexable, index, value object.Object) object.Object {
	hash, isHash := indexable.(*object.Hash)
	array, isArray := indexable.(*object.Array)

	if isHash {
		return evalHashKeyAssignment(hash, index, value)
	}

	if isArray {
		return evalArrayIndexAssignment(array, index, value)
	}

	return newError("Invalid type passed, expected a type of Hash or Array, got %T instead", indexable.Type())
}

func evalHashKeyAssignment(hash *object.Hash, index, value object.Object) object.Object {
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable value as hash key: %s", index.Type())
	}

	hashed_key := key.HashKey()

	hash.Pairs[hashed_key] = object.HashPair{Key: index, Value: value}

	return value
}

func evalArrayIndexAssignment(array *object.Array, index, value object.Object) object.Object {
	idx, ok := index.(*object.Integer)

	if !ok {
		return newError("Invalid index value passed, expected an integer, got: %T", index.Type())
	}

	array.Elements[idx.Value] = value

	return value
}

func applyForLoop(forLoop *ast.ForLoopStatement, env *object.Environment) object.Object {
	var result object.Object

	//Evaluate the first loop condition
	stopLoop := Eval(forLoop.LoopCondition, env)
	loopCondition, ok := stopLoop.(*object.Boolean)

	if !ok {
		return newError("Invalid loop condition type: %T", loopCondition)
	}

	// While this loop condition is false
	for loopCondition.Value {
		result = Eval(forLoop.LoopBlock, env)

		updateVal := Eval(forLoop.CounterUpdate.Value, env)

		env.Set(forLoop.CounterVar.Name.Value, updateVal)

		//Continue evaluating the loop condition in the loop
		loopCondition := Eval(forLoop.LoopCondition, env).(*object.Boolean)

		if !loopCondition.Value {
			return result
		}
	}
	return result
}

func isArray(o object.Object) bool {
	return o.Type() == object.ARRAY_OBJ
}

func isInteger(o object.Object) bool {
	return o.Type() == object.INTEGER_OBJ
}

func isString(o object.Object) bool {
	return o.Type() == object.STRING_OBJ
}

func isHash(o object.Object) bool {
	return o.Type() == object.HASH_OBJ
}

/**
dev notes:

eval()
- should start at the top of the AST, recieving an *ast.Program and
  then traverse every node in it and evaluate every statement (aka tree walking).

self-evaluating expressions:
- what we call literals
- we input an integer into eval() and get that integer back (hence they evaluate themselves.)
- we input an *ast.IntegerLiteral, eval() returns an *object.Literal with a Value of that integer

integer vs boolean comparison:
- with this current implementation boolean comparison will always be faster than integer comparison
- this is because with boolean comparisons we always use pointers to the two boolean object (TRUE, FALSE)
- whereas with integers a new object.Integer has to be instantiated, creating new pointers
- we cannot compare these pointers to different object.Integer instances, otherwise 7==7 would be false
- this is also why integer operands have to be higher up in the switch statement check (see evalInfixExpression)
- as long as we take care of the operand types before arriving at pointer comparisons this evaluation will work fine

internal error handling:
- in order to prevent errors from being passed around and bubbling up from their origin we check for errors whenever we call Eval inside of Eval
**/
