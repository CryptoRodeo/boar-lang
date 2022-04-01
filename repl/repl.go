package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

const BANANA = `
                ██                                      
              ██  ██                                    
            ██  ░░░░██                                  
          ██  ░░░░░░░░██                                
          ██  ░░░░░░░░██                                
          ██  ██░░██░░██                                
          ██  ▓▓░░██░░██                                
          ██  ██░░██░░██                                
          ██  ░░░░░░░░██                                
  ██    ████  ░░░░░░░░████    ██                        
██░░██████  ░░░░██░░░░░░██████░░██                      
██░░░░░░░░░░░░░░██░░░░░░░░░░░░░░██                      
  ██░░░░░░░░░░██████░░░░░░░░░░██                        
    ▓▓▓▓▓▓▓▓▓▓      ▓▓██▓▓▓▓▓▓                          

`

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
		// The line we just read
		line := scanner.Text()
		// pass it through the lexer
		l := lexer.New(line)
		// pass lexed line through the parser
		p := parser.New(l)
		// parse the program
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		//print the currently parsed program
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, BANANA)
	io.WriteString(out, "Whoops, Parser slipped on some errors!\n")
	io.WriteString(out, "Errors found:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
