# monke-lang
*"Return To Monke 🐒"*

Go-based language interpreter for a toy programming language called "monke"

Features:

**Basic math operations:**
```
~> 5 + 5
10

~> 8 * 2
16

~> (1 > 2) == false
true

~> !true
false

```

**Conditional expressions:**
```
~> if (1 > 2) { "a" } else { "b" }
b
```

**Evaluating boolean expresisons:**
```
~> 1 < 2
true
~> true == false
false
```

**dynamic typing:**
```
~> let x = 2
~> x
2

~> let x = "hello"
~> x
hello
```

**string concatenation:**
```
~> "Hello" + " " + "World"
Hello World
```

**functions:**
```
~> let adder = fn(a,b) { a + b }
~> adder(2,3)
5
~> fn(x) { x * 2; }(2)
4
```

**closures**
```
~> let newAdder = fn(x) { fn(y) { x + y } };
~> let addTwo = newAdder(2);
~> addTwo(3);
5
~> let addThree = newAdder(3);
~> a
```

**first class functions**
```
~> let add = fn(a, b) { a + b };
~> let sub = fn(a, b) { a - b };
~> let applyFunc = fn(a, b, func) { func(a, b) };
~> applyFunc(2, 2, add);
4
~> applyFunc(10, 2, sub);
8
```

**Arrays and built in functions:**
```
~> let x = [1,2,3]
~> x
[1, 2, 3]
~> let y = push(x,4)
~> y
[1, 2, 3, 4]
~> len(y)
4
~> len(x)
3
~> first(x)
1
~> last(y)
4
~> x[2]
3
```

**Hash Maps:**
```
~> let person = { "name": "John", "age": (2*15) }
~> person
{age: 30, name: John}
~> person["name"]
John
~> person["age"]
30
```
