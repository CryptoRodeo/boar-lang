package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"

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
	repl.Start(os.Stdin, os.Stdout)
}
