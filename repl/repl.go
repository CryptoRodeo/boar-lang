package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/setuphelpers"
	"os/user"

	"github.com/TwiN/go-color"
)

const PROMPT = "~> "

const TERMINATOR = "exit()"

// Creates a new scanner, object environment and preloads
// built in functions into the global environment
func setup(in io.Reader, out io.Writer) (*bufio.Scanner, *object.Environment) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	setuphelpers.LoadBuiltInMethods(env)

	return scanner, env
}

func Start(in io.Reader, out io.Writer) {
	printInterpreterPrompt()
	scanner, env := setup(in, out)
	// Loop forever, until we exit
	for {
		fmt.Printf("%s", PROMPT)
		scanned := scanner.Scan()

		// If a new line is encountered, don't do anything.
		if !scanned {
			return
		}
		// Grab the line we just read
		line := scanner.Text()

		// Exit
		if line == TERMINATOR {
			break
		}

		// pass it through the lexer
		l := lexer.New(line)
		// pass lexer generated tokens to the parser
		p := parser.New(l)
		// parse the program
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			setuphelpers.PrintParserErrors(out, p.Errors())
			continue
		}

		//print the currently evaluated program
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			// apply syntax highlighting
			str := setuphelpers.ApplyColorToText(evaluated.Inspect())
			io.WriteString(out, str)
			io.WriteString(out, "\n")
		}

	}
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
