# Boar lang 🐗

```
                  #         
          # # # # # #       
      # # #         #       
    # #             # # #   
    #                 # # # 
    # #   # # #   # # # #   
    # # # # # # # # #       
    # # #     # # # #       
```

Language interpreter for a toy programming, built in Go

![workflow](https://github.com/CryptoRodeo/boar-lang/actions/workflows/workflow.yml/badge.svg)

## Details
Based on ["Writing An Interpreter In Go" by Thorsten Ball](https://interpreterbook.com/) with some extra improvements, such as:
- The ability to read and evaluate `.br` code files (or trigger the REPL using `--prompt`)
- Additional built in functions for the Hash and Array objects (inspired from other languages such as Ruby)
- Standard Object#Function invocation: `someObject.someMethod()` as opposed to `someMethod(someObject)`
- Variable reassignment (`let x = 3; x = "hello"` as opposed to `let x = 3; let x = "hello"`)
- Index reassignment for Arrays and Hashes (`hash[key] = expression`, `arr[index] = expression`)
- For loops
- Improved REPL: 
  - evaluate multiple lines
  - user input history
  - syntax highlighting
  - exit typing `exit()`
- Base project refactors
- Additional dev notes for each interpreter component

## Quick Start Guide:

The recommended way is to use Docker:
```
git clone git@github.com:CryptoRodeo/boar-lang.git

cd ./boar-lang

docker build . -t boar-lang

docker run -it --name="boar-lang" boar-lang

# To start the prompt type './boar --prompt'
$ ./boar --prompt
Hello boar, (type 'exit()' to exit)
~> 

# running an .br file (a test file exists)
$ ./boar -f ./test.br

```

You can also just run it regularly (requires go version >= 1.16):
```
# Build executable
go build -o boar

# Running the prompt
$ ./boar --prompt

Hello kilgore, (type 'exit()' to exit)
~> 

# Running a .br file (a test file exists)
$ ./boar -f ./test.br

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

~> x = "hello"
~> x
hello
```

**string concatenation:**
```
~> "Hello" + " " + "World"
Hello World
```

**Error handling:**
```
~> let x

🐗 Error!:
> expected next token to be =, got EOF instead

~> let arr = [1,2 

🐗 Error!:
> expected next token to be ], got EOF instead

~> 
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
~> addThree(5)
8
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

**For loops**
```
~> let y = 0;
~> for (let x = 0; x < 10; x = x + 1) { y = x; };
~> y
9

~> for (let a = 0; a < 5; a = a + 1) { puts(a); }
0
1
2
3
4
```

**Arrays:**
```
::Creating an array
~> let x = [1,2,3]
~> x
[1, 2, 3]

::Adding to the array
~> let y = push(x,4)
~> y
[1, 2, 3, 4]

::Array#len
~> len(y)
4
~> len(x)
3
::Array#len alternative
~> y.len()
4

::Array#first
~> x.first()
1

::Array#last
~> y.last()
4

::Array#[]
~> x[2]
3

::Array index assignment
~> x[2] = "Hello!"
Hello!
~> x
[1, 2, Hello!]

::Array#map
~> let arr = [1,2,3]
~> let addTwo = fn(x) { x + 2; }
~> let res = arr.map(addTwo)
~> res
[3, 4, 5]

::Array#pop
~> let arr = [1,2,3]
~> let popVal = arr.pop()
~> arr
[1, 2]
~> popVal
3

::Array#shift
~> let tb = ["Tom", "Bombadil"]
~> let firstName = tb.shift()
~> firstName
Tom
~> tb
[Bombadil]

::Array#slice
~> let animals = ["ant", "bison", "camel", "duck", "elephant"];
~> animals.slice(2)
[camel, duck, elephant]
~> animals.slice(2, 4)
[camel, duck]
~> animals.slice()
[ant, bison, camel, duck, elephant]
```

**Hash Maps:**
```
::Creating a hash
~> let person = { "name": "John", "age": (2*15) }
~> person
{age: 30, name: John}

::Hash#[]
~> person["name"]
John
~> person["age"]
30

::Hash Index/Key assignment
~> let USDrinkingAge = 21
~> person["age"] = USDrinkingAge
21
~> person["age"]
21

::Hash#valuesAt
~> person.valuesAt("age", "name")
[21, John]

::Hash#toArray
~> person.toArray()
[name, John, age, 21]

::Hash#delete
~> person.delete("age")
{name: John, null: null}
~> person["age"]
null

::Hash#dig
~> let person = { "name": "Tom Bombadil", "clothes": { "shoes": "yellow boots" } };
~> person.dig("clothes", "shoes")
yellow boots

```
## Implementation Details:
- This interpreter uses a tree-walking strategy, starting at the top of the AST, traversing every AST Node and then evaluating its statement(s)
- The parser uses the Vaughan Pratt parsing implementation of associating parsing functions with different token types as well as handling different precedence levels.

