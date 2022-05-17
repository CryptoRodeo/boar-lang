# monke-lang
*"Return To Monke 🐒"*

Go-based language interpreter for a toy programming language called "monke"

## Implementation Details:
- This interpreter uses a tree-walking strategy, starting at the top of the AST, traversing every AST Node and then evaluating its statement(s)
- The parser uses the Vaughan Pratt parsing implementation of associating parsing functions with different token types as well as handling different precedence levels.

## How to run

The recommended way is to use Docker:
```
docker build . -t monke-lang
docker run -it monke-lang --name="monke-lang"

Hello root, feel free to type in commands
~> 
```

You can also just run it regularly:
```
go run .

Hello kilgore, feel free to type in commands
~> 
```

## Language Features:

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
And **much more**

Feel free to explore the code base!
