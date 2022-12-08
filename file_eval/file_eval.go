package file_eval

import (
	"boar/evaluator"
	"boar/lexer"
	"boar/object"
	"boar/parser"
	"boar/setuphelpers"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func EvaluateFile(in io.Reader, out io.Writer, filePath string) {
	env := object.NewEnvironment()
	setuphelpers.LoadBuiltInMethods(env)

	fileContent := locateFile(filePath)
	// pass it through the lexer
	l := lexer.New(fileContent)
	// pass lexer generated tokens to the parser
	p := parser.New(l)
	// parse the program
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		setuphelpers.PrintParserErrors(out, p.Errors())
		return
	}

	//print the currently evaluated program
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		// apply syntax highlighting
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}

}

func locateFile(filePath string) string {
	path := formatUserFilePathInput(filePath)
	return findFile(path)
}

func formatUserFilePathInput(filePath string) string {
	var fullFilePath string
	pwd, _ := os.Getwd()

	// ./fileName.br
	if filePath[0] == '.' && filePath[1] == '/' {
		fullFilePath = filePath
	} else {
		fullFilePath = pwd + "/" + filePath
	}

	fileExtension, validFile := validateFileExtension(fullFilePath)

	if !validFile {
		log.Fatalf("Invalid file type passed, expected a .br file, got %s instead", fileExtension)
	}

	return fullFilePath
}

func findFile(filePath string) string {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}

func validateFileExtension(fPath string) (string, bool) {
	fileExtension := filepath.Ext(fPath)
	return fileExtension, fileExtension == ".br"
}
