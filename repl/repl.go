package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = "~> "

const BEAR = `ʕ•ᴥ•ʔ`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, BEAR)
	io.WriteString(out, "Hello friend, something went wrong\n")
	io.WriteString(out, "Errors found:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
