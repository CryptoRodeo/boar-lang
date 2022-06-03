package main

import (
	"bytes"
	"fmt"
	"monkey/file_eval"
	"monkey/repl"
	"os"

	prompt "github.com/c-bata/go-prompt"
)

func executor(t string) {
	if t == "exit()" {
		fmt.Printf("Goodbye!")
		os.Exit(0)
	}
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func main() {
	// no arguments passed
	if len(os.Args) == 1 {
		printHelpMenu()
		return
	}

	switch os.Args[1] {
	case "--prompt":
		repl.Start()
	case "-f":
		file_eval.EvaluateFile(os.Stdin, os.Stdout, os.Args[2])
	default:
		printHelpMenu()
	}
}

func printHelpMenu() {
	var out bytes.Buffer
	out.WriteString("--prompt to use the interpreter\n")
	out.WriteString("-f FILE to evaluate a .mk file\n")
	fmt.Println(out.String())
}
