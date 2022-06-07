package repl

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/setuphelpers"
	"os"
	"os/user"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/c-bata/go-prompt"
)

// Used to determine if we should swap out the cursor / prefix
// i.e. when evaluating the next line
var LivePrefixState struct {
	LivePrefix string
	IsEnabled  bool
}

var CURSOR = "~> "
var BLOCK_CURSOR = "... "

const TERMINATOR = "exit()"

// Global obj.Environment. Holds builtin functions
var ENV = setupEnv()

// Holds all user input lines, used in case we need to evaluate user input
// on the next line.
var CODE_BUFFER = []string{}

// used to determine if we should evaluate the next line
var CHARS_STILL_OPEN int = 0

func Start() {
	printInterpreterPrompt()

	cursor := prompt.OptionPrefix(CURSOR)
	liveCursor := prompt.OptionLivePrefix(changeLivePrefix)

	p := prompt.New(readInput, completer, cursor, liveCursor)
	p.Run()
}

func shouldContinue(char rune) bool {
	if char == '{' || char == '(' {
		CHARS_STILL_OPEN++
	}

	if char == '}' || char == ')' {
		CHARS_STILL_OPEN--
	}

	return CHARS_STILL_OPEN > 0
}

// Creates a new object environment and preloads
// built in functions into the global environment
func setupEnv() *object.Environment {
	env := object.NewEnvironment()
	setuphelpers.LoadBuiltInMethods(env)
	return env
}

func readInput(line string) {
	if line == "exit()" {
		exitRepl()
	}
	evaluate(line)
}

func completer(t prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "let", Description: "declare a statement"},
		{Text: "puts", Description: "print a value"},
		{Text: "fn", Description: "declare a function literal"},
		{Text: "if", Description: "declare a conditional statement"},
	}

	// Check if we're evaluating the last block, reset cursor so indentation is correct.
	finalChar := getFinalChar(t.CurrentLine())
	if finalBlock(finalChar) {
		resetCursor()
		resetBlockCounter()
	}

	return prompt.FilterHasPrefix(s, t.CurrentLine(), true)
}

func printInterpreterPrompt() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	terminator := color.Ize(color.Red, "exit()")
	userName := color.Ize(color.Cyan, user.Username)
	fmt.Printf("Hello %s, (type '%s' to exit)\n", userName, terminator)
}

func printParserErrors(errors []string) {
	fmt.Print("\n" + setuphelpers.MONKE + " Error!:\n")
	for _, msg := range errors {
		fmt.Print("> " + msg + "\n\n")
		fmt.Println()
	}
}

func evaluate(line string) {
	CODE_BUFFER = append(CODE_BUFFER, line)

	if shouldContinue(getFinalChar(line)) {
		setBlockCursor()
		return
	}

	resetCursor()

	code := formatLine(CODE_BUFFER)
	emptyCodeBuffer()
	// pass it through the lexer
	l := lexer.New(code)
	// pass lexer generated tokens to the parser
	p := parser.New(l)
	// parse the program
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	//print the currently evaluated program
	evaluated := evaluator.Eval(program, ENV)
	if evaluated != nil {
		// apply syntax highlighting
		str := setuphelpers.ApplyColorToText(evaluated.Inspect())
		fmt.Println(str)
	}
}

func emptyCodeBuffer() {
	CODE_BUFFER = make([]string, 0)
}

func formatLine(lines []string) string {
	return strings.Join(lines, " ")
}

func getFinalChar(line string) rune {
	if len(line) == 0 {
		return 0
	}
	return rune(line[len(line)-1])

}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnabled
}

func finalBlock(r rune) bool {
	return r == '}' && CHARS_STILL_OPEN == 1
}

func resetCursor() {
	LivePrefixState.IsEnabled = false
	LivePrefixState.LivePrefix = CURSOR
}

func setBlockCursor() {
	indentationLevel := ""
	for i := 0; i < CHARS_STILL_OPEN; i++ {
		indentationLevel += " "
	}
	LivePrefixState.IsEnabled = true
	LivePrefixState.LivePrefix = (indentationLevel + BLOCK_CURSOR)
}

func resetBlockCounter() {
	CHARS_STILL_OPEN = 0
}

func exitRepl() {
	fmt.Printf("Goodbye!")
	os.Exit(0)
}
