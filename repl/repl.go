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

	"github.com/TwiN/go-color"
	"github.com/c-bata/go-prompt"
)

const PROMPT = "~> "

const TERMINATOR = "exit()"

var ENV = setupEnv()

func Start() {
	printInterpreterPrompt()
	setCursor := prompt.OptionPrefix(PROMPT)
	p := prompt.New(readInput, completer, setCursor)
	p.Run()
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
	// pass it through the lexer
	l := lexer.New(line)
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

func exitRepl() {
	fmt.Printf("Goodbye!")
	os.Exit(0)
}
