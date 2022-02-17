package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

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
		// print all the tokens the lexer gives us until we encounter EOF.
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
