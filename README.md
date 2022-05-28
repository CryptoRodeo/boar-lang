# monke-lang
> "[Reject humanity, return to monke](https://knowyourmeme.com/memes/return-to-monke)" 🐒

Go-based language interpreter for a toy programming language called "monke" (pronounced "monk")

Based on ["Writing An Interpreter In Go" by Thorsten Ball](https://interpreterbook.com/) with some extra improvements, such as:
- Additional built in functions for the Hash and Array objects (inspired from other languages such as Ruby)
- Index reassignment for Arrays and Hashes
- Base project refactors
- Additional dev notes for each interpreter component
- Other features that are currently a WIP

## Quick Start Guide:

The recommended way is to use Docker:
```
git clone git@github.com:CryptoRodeo/monke-lang.git

cd ./monke-lang

docker build . -t monke-lang

docker run -it monke-lang --name="monke-lang"

Hello monke, (use Ctrl+C to exit)
~> 

```

You can also just run it regularly (requires go version >= 1.16):
```
go run .

Hello kilgore, (use Ctrl+C to exit)
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

::Array#first
~> first(x)
1

::Array#last
~> last(y)
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
~> let res = map(arr, addTwo)
~> res
[3, 4, 5]

::Array#pop
~> let arr = [1,2,3]
~> let popVal = pop(arr)
~> arr
[1, 2]
~> popVal
3
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
~> valuesAt(person, "age", "name")
[21, John]

::Hash#toArray
~> toArray(person)
[name, John, age, 21]

::Hash#delete
~> delete(person, "age")
{name: John, null: null}
~> person["age"]
null

::Hash#dig
~> let person = { "name": "Tom Bombadil", "clothes": { "shoes": "yellow boots" } };
~> dig(person, "clothes", "shoes")
yellow boots

```
## Implementation Details:
- This interpreter uses a tree-walking strategy, starting at the top of the AST, traversing every AST Node and then evaluating its statement(s)
- The parser uses the Vaughan Pratt parsing implementation of associating parsing functions with different token types as well as handling different precedence levels.

