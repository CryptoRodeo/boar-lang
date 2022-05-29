package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
	"strings"

	"github.com/TwiN/go-color"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	ctrlC := color.Ize(color.Red, "Ctrl+C")
	terminator := color.Ize(color.Red, "exit()")
	userName := color.Ize(color.Cyan, user.Username)
	fmt.Printf("Hello %s, use (%s or type '%s' to exit)\n", userName, ctrlC, terminator)

	switch os.Args[1] {

	case "--prompt":
		repl.Start(os.Stdin, os.Stdout)
	default:
		EvaluateFile(os.Stdin, os.Stdout, os.Args[1])
	}
}

func EvaluateFile(in io.Reader, out io.Writer, fileName string) {
	env := object.NewEnvironment()
	repl.LoadBuiltInMethods(env)
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range strings.Split(string(content), "\n") {
		// pass it through the lexer
		l := lexer.New(code)
		// pass lexer generated tokens to the parser
		p := parser.New(l)
		// parse the program
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(out, p.Errors())
			continue
		}

		//print the currently evaluated program
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			// apply syntax highlighting
			str := repl.ApplyColorToText(evaluated.Inspect())
			io.WriteString(out, str)
			io.WriteString(out, "\n")
		}

	}
}
