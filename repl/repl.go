package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os/user"
)

const PROMPT = "~> "

const BEAR = `ʕ•ᴥ•ʔ`

const TERMINATOR = "exit()"

// Creates a new scanner, object environment and preloads
// built in functions into the global environment
func setup(in io.Reader, out io.Writer) (*bufio.Scanner, *object.Environment) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	loadBuiltInMethods(env)

	return scanner, env
}

func loadBuiltInMethods(env *object.Environment) {
	for key, value := range evaluator.BUILTIN {
		env.Set(key, value)
	}
}

func Start(in io.Reader, out io.Writer) {
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
			printParserErrors(out, p.Errors())
			continue
		}

		//print the currently evaluated program
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	io.WriteString(out, "\n"+BEAR+"\n")
	io.WriteString(out, "psst, hey "+user.Username+", I think you broke something...\n")
	io.WriteString(out, "Errors found:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
