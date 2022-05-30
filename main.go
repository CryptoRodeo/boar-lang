package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"monkey/setuphelpers"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
)

func main() {
	// no arguments passed
	if len(os.Args) == 1 {
		printHelpMenu()
		return
	}

	switch os.Args[1] {
	case "--prompt":
		printInterpreterPrompt()
		repl.Start(os.Stdin, os.Stdout)
	case "-f":
		evaluateFile(os.Stdin, os.Stdout, os.Args[2])
	default:
		printHelpMenu()
	}
}

func evaluateFile(in io.Reader, out io.Writer, filePath string) {
	env := object.NewEnvironment()
	setuphelpers.LoadBuiltInMethods(env)

	fileContent := locateFile(filePath)
	for _, code := range fileContent {
		// pass it through the lexer
		l := lexer.New(code)
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

func locateFile(filePath string) []string {
	return formatUserFilePathInput(filePath)
}

func formatUserFilePathInput(filePath string) []string {
	var errorFound error
	var fileContents []byte
	var fullFilePath string
	pwd, _ := os.Getwd()

	// ./fileName.mke
	if filePath[0] == '.' && filePath[1] == '/' {
		fullFilePath = filePath
	} else {
		fullFilePath = pwd + "/" + filePath
	}

	fileExtension, validFile := validateFileExtension(fullFilePath)

	if !validFile {
		log.Fatalf("Invalid file type passed, expected a .mke file, got %s instead", fileExtension)
	}

	fileContents, errorFound = ioutil.ReadFile(fullFilePath)
	if errorFound != nil {
		log.Fatal(errorFound)
	}

	// Split by ';'
	res := strings.Split(string(fileContents), ";")

	return res
}

func validateFileExtension(fPath string) (string, bool) {
	fileExtension := filepath.Ext(fPath)
	return fileExtension, fileExtension == ".mke"
}
func printHelpMenu() {
	var out bytes.Buffer
	out.WriteString("--prompt to use the interpreter\n")
	out.WriteString("-f FILE to evaluate a .mke file\n")
	fmt.Println(out.String())
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
