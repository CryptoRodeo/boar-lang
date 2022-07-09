43a6e00 (HEAD -> main, origin/main, origin/HEAD) Used wrong single quote type...
9958f33 Update JenkinsFile with an additional trigger, have it poll SCM for any changes
030d259 Small updates
194b4e1 (tag: v4.1) add demo gif, update readme
1de1957 evaluate last char to determine if we're in a block statement, dynamically add indentation when evaluating multiple lines
ae824f5 (tag: v4.0) update README.md
4e718f9 Merge pull request #15 from CryptoRodeo/repl-improvements
1b68189 remove unused methods
fa772f0 add some comments, change prefix live
3e40dad update README
d4df03d It ugly, but now we an evaluate multiple lines (WIP)
bf85235 Heavy WIP - Upgrade REPL, use GoPrompt package
80faa24 (tag: v3.5) Merge pull request #14 from CryptoRodeo/for-loops
2d7c09a move intepreter prompt func to where its actually used
734990e update file extensions
978d602 we'll evaluate '.mk' files now
e1a7806 create file_eval package, evlauate files here
472ef92 refactor main file
3bb0978 update README.md
e5c4e29 evaluate for loops
d498b98 fix test output
8c857e1 add test to evaluate for loops
66eaa53 ehh, good enough test.
dc7e376 Update README.md
a86de83 Update README.md
2a65809 clean up parser
1be154b D'oh!
99783db WIP: parse for loops
2fbec46 WIP: test to parse for loops
46bc3ce update For loop printing
adf6b8e change this out to LoopBlock
c419cd8 (origin/for-loops, for-loops) WIP: add method to parse for loops
f880dba D'oh!
258bfc4 add for loop AST node
30ebfae test tokenizing / lexing a for loop
8ac162d add keyword 'for' to token list
8dca00f update files
e25ecae move coloring functions to setuphelpers module
5467a5c delete this
bd4de47 update gitignore
de2c1e2 update test.mke
a9de2d7 (tag: v3.0) Merge pull request #13 from CryptoRodeo/evaluate-code-files
c9eee61 updated README.md
f2bdeb8 build the executable in the container
b53261a make the container entrypoint bash
576d8d1 (origin/evaluate-code-files, evaluate-code-files) update README.md
8ad43af Update test.mke file
9a3f404 test.mke file in test directory to test evaluating file paths
247f87a test.wrong file ("To check file extension validation")
575b788 test .mke file
3ec281b refactor
4072869 WIP update README.MD, needs quick start guide update
09f6835 refactor main to handle interpreter prompt and reading from files
26fe02d move these functions to this helper module
975dae0 WIP, some of these functions will have to be moved out since they'll also be used when evaluating files.
3ff809d test mke file
8c75b14 Add ability to read from .mke files
05890af Update README.md
134a1bb Update CHANGELOG.md
4002f5e (tag: v2.5) Merge pull request #12 from CryptoRodeo/variable-reassignment
d0157a4 (origin/variable-reassignment, variable-reassignment) update README.md
89d77e9 evaluate assignment expressions
8bc2764 add test to evaluate assignment expressions
2fdbd92 - assign precedence value to ASSIGN token - register ASSIGN token prefix - register infix parse function for ASSIGN token
42892a4 add test to parse assignment statements
2ccef3a add lexer test to tokenize assignment statements
c4b63a1 fix minor blunder, add AssignmentExpression node
7f96e4d tidy up go.mod
6765cb1 update comment
e9d7414 remove unused struct value
4b87d2a Merge pull request #11 from CryptoRodeo/improve-prompt
c60320b update error handling format
9cea3a4 update examples, add error handling examples
4a07b97 update Error message format
8ec5432 update README.md
2d52cd9 remove history arr. Removing the history traversal idea from this PR.
ab1cf22 move syntaxHighlighter to REPL package since thats where its used.
7f3ed50 update README.md
d69014a update how hashes are printed
5d7f88f apply syntax highlighting in REPL
1a0a4c1 create helper module to highlight syntax in REPL
f77d2ca add color package to main file, colorize prompt text
2c56209 add imported modules
450f20e update README
def5469 update start prompt message
5c7de3c add REPL loop terminator
b16049f Merge pull request #10 from CryptoRodeo/function-call-improvements
99f643c update README to use new built in function invocation sytax
0cee69c add additional test scenarios for the new function call syntax
171e8ea implement evaluating internal function calls
2cacadb add test for array internal function call
f5e6719 WIP: Begin evaluating internal function calls
acd88db update example comments
31677bd - create a new INTERNAL_CALL precedence value, set to highest. - register an infix parsing function for the dot token - Create method for parsing interal function calls
cd7c59d create test for parsing internal function calls
fa2c7fe alter InternalFunctinCall AST Node
81665a5 WIP add InternalFunctionCall AST Node type
c3a3649 tokenize the '.' character
9fdd8d6 Move this to the delimeter category
6a911f9 create a test to tokenize the '.' char in the lexer
75b70d8 actually lets rename this to dot. no one says "array period slice()"
3d62602 add '.' character to token const dictionary
357c6e0 (tag: 2.0) Merge pull request #9 from CryptoRodeo/array-improvements
3ed7512 add Array#slice example to README.md
123fa0b Implement Array#slice
7dcff14 remove inspect line
dbe8be2 whoops, forgot to add this
2594d84 add Array#shift to README.md
d8d25d4 implement Array#shift
385f3ec add tests for Array#shift
81fbf4d Give Credit Where Credit Is Due :)
744dacf update spacing in README
f7322cb update README.md
8f5f9ac working Array#pop built in method
89f9ff7 update test for Array#pop
d46bb0a add test for Array#pop
43ad395 more comments
0183d90 add comments for setup function
c5deb17 update README.md
bcec46a WIP: working map built in function!
8e32008 WIP: remove built in method check.  add check in applyFunction for map calls
ad40a9b remove uneeded print statement
923174c preload built in methods in the test environment as well
fc3d743 preload built in methods before starting repl to get around invalid cycle issue with built in methods.
cd7d597 (origin/array-improvements, array-improvements) remove redundant line in README.md
ee40307 add test for map built in function
032d2f6 WIP - map built in function
0930329 update README wiith example of array index assignment
e35d6fd Update start prompt message
7571680 Evaluate array index assignments
3ce00fd create test to evaluate array index assignments
42516f1 add todo comment for first array object improvement
f32db88 Fix FuncName values for error formatter
366d483 fix README.md
8ff4515 (tag: v1.5) update changelog
0de354f refactor slice creation in __dig__
c5080a2 Remove redundent structs, use new hash builtin method error checker
439ae0a prepare hash builtin method error checker
e4a66ec add CHANGELOG file
3874f46 (evaluation) fix README.md example
c729bf1 Merge pull request #8 from CryptoRodeo/hash-improvements
09297d0 update README.md
99ca0c3 WIP: implement the Hash#dig method
833904b add test to dig hashes by keys
60c6927 fix and update README.md
b31b430 Merge branch 'main' into hash-improvements
536aa5a delete this
c312e70 add built in function to convert hash to array
b41d998 add test to convert hash object to array
0fff0f0 finished test for hash value retrieval
e7c692d WIP: add function and test to extract hash values
27daaa1 add new built in function for hash key deletions
dd501aa add test for hash key deletions
1df40e7 remove unused env argument
09b5e41 improve docker file (dont run root kids)
fe60693 Merge branch 'main' of github.com:CryptoRodeo/monke-lang into main
e237161 Update README.md
4411312 small refactor for the index assignment parsing
c40c99e add case in evaluator for hash assignments
fad7395 WIP - Add test to evaluate hash assignments
c3a56ff add IndexAssignment AST Node (WIP)
4541286 add lexer test for parsing hash index assignments
2b677b2 add test for parsing index assignments (WIP)
62b0c89 add case where we might be assigning a value at an index
375c644 Update README.md
dda4429 Update README.md
38c7568 (tag: v1.0) update jenkins script
472d919 fix hash inspect output
8011a61 add hash indexing functionality
f2fed76 add test for hash indexing
172d313 update README.md
7b7b199 Update README.md
bf5ed90 Update README.md
0423f14 Merge pull request #7 from CryptoRodeo/extend
dc23ed7 (origin/extend) finally test puts implementation
b399ad7 implement "puts"
501bc97 update test, still havent implemented puts yet
44f1fd3 refactor, fix failing tests
d7cec8d Evaluate hashes
c3b668f Add test for evaluating hash literals
7088aa6 Whoops
e3bcf6e add hashable interface
a8032a1 Add inspect method for hashes
18125c8 Define hash pair, hash and related methods
bfa4895 add methods to generate hash keys
90fb3bd add test to compare hash keys
0b74f14 add support for parsing hash literals
8bb8982 add tests for parsing hashes
1c157bd add hash literal AST node
9af30c5 lex colons
b80964e add test to parse hash literals
6f99398 Add semicolon token
9ad5f52 add support for 'push' method
50b2b37 add support for 'rest' built in function
ee3e8c5 add tests for built in functions
4269d1b small refactor, add support for 'last' inbuilt function
d9e4096 add support for arrays in the built in 'len' function
f8cb692 small changes
54f7da4 evaluate index expressions + new helper functions
31aeb23 test evaluating index expressions
8c8f365 evaluate arrays
67fb087 test evaluating arrays
0e61fe7 add array object struct + functions
db43b24 parse index expressions
c741aef add test for index expression parsing
393778e lets support index expressions
acb68d5 more comments
8d2e7c1 parse array literals, refactor parsing call expressions
3856e51 define array ast node
5dca914 test parsing array literals
37d6aa6 support creating bracket tokens in the lexer
d0f5d94 test support for array brackets in lexer
c5486e9 let start supporting arrays
7718e27 allow builtin functions to be handled by evaluator
1a641a3 update test
0fe88c8 fully functional len() builtin method
6efcd4e first built-in function: len()
6979f44 lets create some built-in functions
ecdcdac support string concatenation
9b4dc46 add test + regression test for string concat
b04a332 add test case for evaluating string literals
1f4d506 evaluate strings properly
e49caee add strings to the object system
62b980d create & register function for string literal parsing
57d5a3e create test for parsing string literals
b7bd5bb add string literal AST node
eb6885f handle strings in lexer
2355884 extend test
7ccb1f9 add string literal TokenType
88d3a3b Merge pull request #6 from CryptoRodeo/evaluation
56686fb (origin/evaluation) spelling fix
acbff73 yeah boi we got closures too
55fc42b hell yeah, we can now define and call functions!
8d15e76 allow for inner scopes
f125f18 WIP: evaluate call expressions
3ad45a0 add test to test function applications
804b04c evalute those functions
3ed8c94 add test to evaluate functions
35d7302 add definition and methods for func objects
9579299 update repl to use new env object
ede7294 update evaluator to use new env object
ad032f3 test evaluating let statements
5f9f7f5 add environment object
d359577 add tests for internal error handling
c5a23f6 add internal error handling
14fa989 Add error struct and type
5cc4470 use new method for BlockStatements to parse return statements
73158b2 add more complex test scenario
875e57b Add return statement evaluation
69880fa add test to evaluate return statements with ints
2dc95d2 add return value object
f595345 evaluate conditionals
ef9aa0e add test and helper function for conditionals
16b7502 add dev notes
118b32f add support for boolean infix expressions
b99ec5c extend tests for boolean operand evaluation
02d68c2 evaluate infix boolean expressions
edaf76b extend test, add boolean infix expressions
61d33df yeah dawg evaluate integer infix expressions
72a00d0 extend tests for infix operators
385d67e small comment
5609ca2 evaluate minus operator prefix expressions
d82afd1 extend test to prepare parsing prefix expressions
298cef8 smol note
2698914 parse bang operator expressions
c333b61 Test evaluating bang operator prefix expressions
93cef2c add NULL object
c826482 optimize evaluating booleans
f1d200e Merge branch 'evaluation' of github.com:CryptoRodeo/monke-lang into evaluation
a815aad smol note change
66c7e9a add boolean evaluation test
5a98474 update REPL message
a83184b evaluate booleans
886273d update message
69c126c lets now put the E in REPL
b67e1e1 remove phallic banana, change prompt to fish shell style
9bbf47e add note on self-evaluating expressions
cfd82e4 update jenkins test script
914d7da add evaluator, evaluate integers from the AST for now.
fe9554e Add evaluator test
cffdd49 tread safely now, we're implementing null
8beedf1 Adding bools
0841d9d Add integer object
04e1926 Adding object package
6b42e7e have container entrypoint run repl
355b8e4 Merge pull request #5 from CryptoRodeo/update-repl
69fbc25 add error message banana cause why not
49d51a6 more of an RPPL than a REPL
abc705e add more notes
8e5a7f5 Merge pull request #4 from CryptoRodeo/call-expressions
b4de374 (origin/call-expressions) ensure call expressions have highest precedence
f8ced44 assign the correct precedence for lparen token
760766a register infix function for lparen, create methods to parse call expressions + arguments
5fe66c2 add test for parsing call expression parameters
5292510 like I said, coding at 5am be hard
6a13582 add test for parsing call expressions
21fe203 whoops (coding at 5am be hard)
07490aa add CallExpression struct
3c4b445 add comment
87c637a update gitignore
e2389fa Merge pull request #3 from CryptoRodeo/function-literals
63975aa (origin/function-literals) add test for function parameter parsing
a476948 register and create function literal parsing function
fb956c5 fix failing test
3548219 fix missing function body error
f956ec6 add also janky entrypoint shell script
90224cb add WIP janky Dockerfile
8cd9947 update Jenkinsfile formatting
b75a20d update script permissions so Jenkins can run them
2c9c396 update Jenkisfile
5a39098 add helper shell scripts
a31b498 update jenkins file
08165c3 still broken...
d4e1bc6 fix pipeline steps
212349c add actual steps to pipeline
9d3cda7 testing local Jenkins server
a0a478f add additional notes
50b6c32 test for function literals WIP
b94afc0 Add AST node to parse function literals
6a36471 Merge pull request #2 from CryptoRodeo/if-expressions
d212847 (origin/if-expressions) More commentary
2d0379b parse else statements (that was easy...)
a5da2c2 add parseBlockStatement method to parse If statements
f92cb88 typo fix
1b9a780 WIP -> add function to parse if expressions
14e9747 small comment change
6c167fb add test for if else expressions
6d90e47 Add test for if expressions
d1eba2e Add if and block statement AST nodes
d914e36 add additional notes on parser concepts
535790f add additional notes on tokens
6aeae93 Add Dev Notes on lexer + tokens
1906656 Merge pull request #1 from CryptoRodeo/grouped-expressions
ffb626a (origin/grouped-expressions) Create and register prefix parsing function for parenthesis
0f4af4f Add additional test scenarios to test grouped expressions
b042886 register prefix parsing functions for boolean values, update parsing functions to use parseExpression
d0fc4a8 Extend test, add function to test boolean literal
cf719bd add test for boolean values. Should fail because of no prefix parse function mapping
50d98bc Add boolean struct so we can start implementing booleans
d5abc37 Clean up some tests
7859f33 Add new helper functions to build more generalized tests. update some tests
cef7f99 add parser tracer
49e12e9 moar notes
abcd56a update some notes.
e2f0efc (origin/pratt-parser-implementation) add tests for precedence parsing
bed0244 add infix parsing to parseExpression
30932b1 create method to parse infix expression, register tokens to use this method
3655e6b Add notes and initialize map for infix parsing functions
ac3bfee Add precedence-helper methods
5b3adc9 Add precedences table / map
12c988a Add infix expression AST node + functions
1ca9374 (WIP) add test for parsing infix expressions
722ada7 Add more notes for the parsePrefixExpression method
d9c6a95 Add functionality to parse expressions with prefixes
eb27fd3 Add prefix expression struct, functions and test
76d7da9 Add integer literal struct, parsing function and test
81e4150 Add functions and tests for parsing identifiers.
de5d840 add parsing functions and maps
ecc10c9 Adding some more notes..
b22653c (origin/parse-expressions) Add sample test for parser
8eb9ff8 Add ExpressionStatement struct and String() methods
688738a Adding gitignore before moving on
cc2a137 Merge branch 'main' of github.com:CryptoRodeo/monke-lang
f680a43 Adding return statement parsing
39f2e8a Update README.md
8f8225c Adding initial files
47bc95d Initial commit
