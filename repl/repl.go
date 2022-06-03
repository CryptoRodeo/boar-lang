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

var CURSOR = "~> "
var CURSOR_OPTION = prompt.OptionPrefix(CURSOR)

const TERMINATOR = "exit()"

var ENV = setupEnv()

var CODE_BUFFER = []string{}
var charsStillOpen int = 0
var PROMPT *prompt.Prompt = prompt.New(readInput, completer, CURSOR_OPTION)

func Start() {
	printInterpreterPrompt()
	PROMPT.Run()
}

func shouldContinue(code string) bool {
	for _, c := range code {
		if c == '{' || c == '(' {
			charsStillOpen++
		}

		if c == '}' || c == ')' {
			charsStillOpen--
		}
	}
	return charsStillOpen > 0
}

// Creates a new scanner, object environment and preloads
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
	}

	return prompt.FilterHasPrefix(s, t.CurrentLine(), true)
}

func printInterpreterPrompt() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	ctrlC := color.Ize(color.Red, "Ctrl+C")
	terminator := color.Ize(color.Red, "exit()")
	userName := color.Ize(color.Cyan, user.Username)
	fmt.Printf("Hello %s, use (%s or type '%s' to exit)\n", userName, ctrlC, terminator)
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

	if shouldContinue(line) {
		return
	}

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

func exitRepl() {
	fmt.Printf("Goodbye!")
	os.Exit(0)
}
