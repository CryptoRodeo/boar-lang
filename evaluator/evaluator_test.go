package evaluator

import (
	"boar/lexer"
	"boar/object"
	"boar/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func loadBuiltInMethods(env *object.Environment) {
	for key, value := range BUILTIN {
		env.Set(key, value)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	loadBuiltInMethods(env)

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got %d, wanted %d", result.Value, expected)
		return false
	}
	return true
}

func testEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got %t wanted %t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got %T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		// when a conditional doesn't evaluate to a value, return NULL/nil
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`, 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Boar"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got %T (%+v)", evaluated, evaluated)
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, got %q expected %q", errObj.Message, tt.expectedMessage)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not a Function, got %T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Params=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);

	addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String, got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value, got %q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String, got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value, got %q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got 2, wanted 1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, nil},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {

		case int:
			testIntegerObject(t, evaluated, int64(expected))

		case nil:
			testNullObject(t, evaluated)

		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)

			if !ok {
				t.Errorf("obj no Array, got %T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong number of elements, wanted %d got %d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedNum := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedNum))
			}

		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not array, got %T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements, got %d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10-9,
		two: 1 + 1,
		"thr" + "ee": 6/2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)

	if !ok {
		t.Fatalf("Eval didn't return Hash, got %T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs, got %d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]

		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashAssignments(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let hash = {"a": 2 }; hash["a"] = 5; hash["a"]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer := tt.expected
		testIntegerObject(t, evaluated, integer)
	}
}

func TestHashKeyDeletions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let hash = {"a": 2 }; delete(hash, "a"); hash["a"]`, nil},
		{`let hash = {"a": 2 }; hash.delete("a"); hash["a"]`, nil},
		{`let hash = {"a": 2, "b": 3 }; delete(hash, "b"); hash["b"]`, nil},
		{`let hash = {"a": 2, "b": 3 }; hash.delete("b"); hash["b"]`, nil},
		{`let hash = {"a": 2, "b":3, "c": 4 }; delete(hash, "a", "b"); hash["c"]`, 4},
		{`let hash = {"a": 2, "b":3, "c": 4 }; hash.delete("a", "b"); hash["c"]`, 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashValueRetrieval(t *testing.T) {
	expected_results := [][]interface{}{
		{2},
		{2, 3},
		{2, nil},
		{nil, nil, nil},
	}
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`let hash = {"a": 2 }; valuesAt(hash, "a");`, expected_results[0]},
		{`let hash = {"a": 2 }; hash.valuesAt("a");`, expected_results[0]},
		{`let hash = {"a": 2, "b": 3 };  valuesAt(hash, "a", "b")`, expected_results[1]},
		{`let hash = {"a": 2, "b": 3 };  hash.valuesAt("a", "b")`, expected_results[1]},
		{`let hash = {"a": 2, "b": 3 }; delete(hash, "b"); valuesAt(hash,"a","b") `, expected_results[2]},
		{`let hash = {"a": 2, "b": 3 }; hash.delete("b"); hash.valuesAt("a","b") `, expected_results[2]},
		{`let hash = {"a": 2, "b":3, "c": 4 }; delete(hash, "a", "b", "x"); valuesAt(hash, "a","b","c")`, expected_results[3]},
		{`let hash = {"a": 2, "b":3, "c": 4 }; hash.delete("a", "b", "x"); hash.valuesAt("a","b","c")`, expected_results[3]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayValues(t, evaluated, tt.expected)
	}
}

func testArrayValues(t *testing.T, evaluated object.Object, expected []interface{}) bool {
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Errorf("object is not an Array, got %T (%+v)", evaluated, evaluated)
		return false
	}

	for _, val := range expected {
		if !contains(result.Elements, val) {
			t.Errorf("Value not found in array: %s", val)
			return false
		}
	}
	return true
}

func contains(array []object.Object, value interface{}) bool {
	for _, val := range array {
		// convert objects to their types so we can correctly compare values
		intVal, _ := value.(int)
		intObject, isInteger := val.(*object.Integer)
		if isInteger && intObject.Value == int64(intVal) {
			return true
		}

		_, isNull := val.(*object.Null)
		if isNull && value == nil {
			return true
		}

		stringVal, _ := value.(string)
		stringObject, isString := val.(*object.String)
		if isString && stringObject.Value == stringVal {
			return true
		}
	}
	return false
}

func TestHashToArrayConversion(t *testing.T) {
	expected_results := [][]interface{}{
		{"a", 2},
		{"a", 2, "b", 3},
	}
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`let hash = {"a": 2 }; toArray(hash)`, expected_results[0]},
		{`let hash = {"a": 2 }; hash.toArray()`, expected_results[0]},
		{`let hash = {"a": 2, "b": 3 };  toArray(hash)`, expected_results[1]},
		{`let hash = {"a": 2, "b": 3 };  hash.toArray()`, expected_results[1]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayValues(t, evaluated, tt.expected)
	}
}

func TestHashDigging(t *testing.T) {
	expected_results := []string{
		"yellow boots",
	}
	tests := []struct {
		input    string
		expected string
	}{
		{`let person = { "name": "Tom Bombadil", "clothes": { "shoes": "yellow boots" } }; dig(person, "clothes", "shoes");`, expected_results[0]},
		{`let person = { "name": "Tom Bombadil", "clothes": { "shoes": "yellow boots" } }; person.dig("clothes", "shoes");`, expected_results[0]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		val, ok := evaluated.(*object.String)

		if !ok {
			t.Errorf("object is not String, got %T (%+v)", val, val)
		}

		if val.Value != tt.expected {
			t.Errorf("Invalid value returned, expected: %s , got: %s", tt.expected, val)
		}
	}
}

func TestArrayIndexAssignments(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let arr = ["a", "b", "c"]; arr[0] = 5; arr[0]`, 5},
		{`let arr = ["tom", "jerry"]; arr[1] = "bombadil"; arr[1]`, "bombadil"},
		{`let arr = ["Smitty Werbenjagermanjensen","He was number", "two"]; arr[2] = "one"; arr[2]`, "one"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		_, isInt := evaluated.(*object.Integer)
		expectedInt, expectedIsInt := tt.expected.(int)

		actualString, isString := evaluated.(*object.String)
		expectedString, expectedIsString := tt.expected.(string)

		if isInt && expectedIsInt {
			testIntegerObject(t, evaluated, int64(expectedInt))
		}

		if isString && expectedIsString {
			if actualString.Value != expectedString {
				t.Errorf("Expected %s, got %s", tt.expected, evaluated)
			}
		}
	}
}

func TestArrayMapFunction(t *testing.T) {
	expectedResults := [][]interface{}{
		{3, 4, 5},
		{4, 8, 12},
	}
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`let arr = [1,2,3]; let addTwo = fn(x) { x + 2; }; let arr = map(arr, addTwo); arr`, expectedResults[0]},
		{`let arr = [1,2,3]; let addTwo = fn(x) { x + 2; }; let arr = arr.map(addTwo); arr`, expectedResults[0]},
		{`let arr = [2,4,6]; let multiplyTwo = fn(x) { x * 2; }; let arr = map(arr, multiplyTwo); arr`, expectedResults[1]},
		{`let arr = [2,4,6]; let multiplyTwo = fn(x) { x * 2; }; let arr = arr.map(multiplyTwo); arr`, expectedResults[1]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, isArray := evaluated.(*object.Array)

		if !isArray {
			t.Errorf("Expected an Array returned, got %T instead", arr.Type())
		}

		for _, val := range tt.expected {
			if !contains(arr.Elements, val) {
				t.Errorf("Invalid value found in array, expected to find: %s", val)
			}
		}
	}
}

func TestArrayPopFunction(t *testing.T) {
	expectedResults := [][]interface{}{
		{1, 2},
		{2, 4},
	}
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`let arr = [1,2,3]; pop(arr); arr`, expectedResults[0]},
		{`let arr = [1,2,3]; arr.pop(); arr`, expectedResults[0]},
		{`let arr = [2,4,6]; pop(arr); arr`, expectedResults[1]},
		{`let arr = [2,4,6]; arr.pop(); arr`, expectedResults[1]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, isArray := evaluated.(*object.Array)

		if !isArray {
			t.Errorf("Expected an Array returned, got %T instead", arr.Type())
		}

		for _, val := range tt.expected {
			if !contains(arr.Elements, val) {
				t.Errorf("Invalid value found in array, expected to find: %s", val)
			}
		}
	}
}

func TestArrayPopFunctionReturn(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let arr = [1,2,3]; pop(arr);`, 3},
		{`let arr = [1,2,3]; arr.pop();`, 3},
		{`let arr = [2,4,6]; pop(arr);`, 6},
		{`let arr = [2,4,6]; arr.pop();`, 6},
		{`let arr = ["Frodo", "Baggins"]; pop(arr);`, "Baggins"},
		{`let arr = ["Frodo", "Baggins"]; arr.pop();`, "Baggins"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, isInt := evaluated.(*object.Integer)
		expectedInt, expectedIsInt := tt.expected.(int)

		if isInt && expectedIsInt {
			testIntegerObject(t, integer, int64(expectedInt))
		}

		stringObj, isString := evaluated.(*object.String)

		expectedString, expectedIsString := tt.expected.(string)

		if isString && expectedIsString {
			if expectedString != stringObj.Value {
				t.Errorf("Invalid string value returned from Array#pop, expected to find: %s, got: %s instead", tt.expected, stringObj.Value)
			}
		}
	}
}

func TestArrayShiftFunction(t *testing.T) {
	expectedResults := [][]interface{}{
		{2, 3},
		{4, 6},
	}
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`let arr = [1,2,3]; shift(arr); arr`, expectedResults[0]},
		{`let arr = [1,2,3]; arr.shift(); arr`, expectedResults[0]},
		{`let arr = [2,4,6]; shift(arr); arr`, expectedResults[1]},
		{`let arr = [2,4,6]; arr.shift(); arr`, expectedResults[1]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, isArray := evaluated.(*object.Array)

		if !isArray {
			t.Errorf("Expected an Array returned, got %T instead", arr.Type())
		}

		for _, val := range tt.expected {
			if !contains(arr.Elements, val) {
				t.Errorf("Invalid value found in array, expected to find: %s", val)
			}
		}
	}
}

func TestArrayShiftFunctionReturn(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let arr = [1,2,3]; shift(arr);`, 1},
		{`let arr = [1,2,3]; arr.shift();`, 1},
		{`let arr = [2,4,6]; shift(arr);`, 2},
		{`let arr = [2,4,6]; arr.shift();`, 2},
		{`let arr = ["Frodo", "Baggins"]; shift(arr);`, "Frodo"},
		{`let arr = ["Frodo", "Baggins"]; arr.shift();`, "Frodo"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, isInt := evaluated.(*object.Integer)
		expectedInt, expectedIsInt := tt.expected.(int)

		if isInt && expectedIsInt {
			testIntegerObject(t, integer, int64(expectedInt))
		}

		stringObj, isString := evaluated.(*object.String)

		expectedString, expectedIsString := tt.expected.(string)

		if isString && expectedIsString {
			if expectedString != stringObj.Value {
				t.Errorf("Invalid string value returned from Array#pop, expected to find: %s, got: %s instead", tt.expected, stringObj.Value)
			}
		}
	}
}

func TestArraySliceFunction(t *testing.T) {
	expectedResults := [][]string{
		{"camel", "duck", "elephant"},
		{"camel", "duck"},
		{"bison", "camel", "duck", "elephant"},
		{"ant", "bison", "camel", "duck", "elephant"},
	}
	tests := []struct {
		input    string
		expected []string
	}{
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; slice(animals, 2);`, expectedResults[0]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; slice(animals, 2, 4);`, expectedResults[1]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; slice(animals, 1, 5);`, expectedResults[2]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; slice(animals);`, expectedResults[3]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, isArray := evaluated.(*object.Array)

		if !isArray {
			t.Errorf("Expected an Array returned, got %T instead", arr.Type())
		}

		for _, val := range tt.expected {
			if !contains(arr.Elements, val) {
				t.Errorf("Invalid value found in array, expected to find: %s", val)
			}
		}
	}
}

func TestArrayInternalFunctionCall(t *testing.T) {
	expectedResults := [][]string{
		{"camel", "duck", "elephant"},
		{"camel", "duck"},
		{"bison", "camel", "duck", "elephant"},
		{"ant", "bison", "camel", "duck", "elephant"},
	}
	tests := []struct {
		input    string
		expected []string
	}{
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; animals.slice(2);`, expectedResults[0]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; animals.slice(2, 4);`, expectedResults[1]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; animals.slice(1, 5);`, expectedResults[2]},
		{`let animals = ["ant", "bison", "camel", "duck", "elephant"]; animals.slice();`, expectedResults[3]},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, isArray := evaluated.(*object.Array)

		if !isArray {
			t.Errorf("Expected an Array returned, got %T instead", arr.Type())
		}

		for _, val := range tt.expected {
			if !contains(arr.Elements, val) {
				t.Errorf("Invalid value found in array, expected to find: %s", val)
			}
		}
	}
}
func TestAssignmentExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a = 4; a;", 4},
		{"let a = 5 * 5; a = 4 * 4; a;", 16},
		{"let a = 5; let b = a; b = 300; b;", 300},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestForLoopStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"for(let x = 0; x < 10; x = x + 1) { puts(x) };", nil},
		{"let y = 0; for(let x = 0; x < 10; x = x + 1) { y = x; }; y;", 9},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if tt.expected == nil {
			nullObj, isNull := evaluated.(*object.Null)
			if !isNull {
				t.Errorf("Expected to get back a null value, got %T back instead", nullObj)
			}
		}

		if intVal, ok := tt.expected.(int); ok {
			intObj, isInt := evaluated.(*object.Integer)

			if !isInt {
				t.Errorf("Expected to get back an integer value, got %T back instead", intObj)
			}
			testIntegerObject(t, intObj, int64(intVal))
		}

	}
}
