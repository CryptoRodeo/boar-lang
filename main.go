package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s, (use Ctrl+C or type 'exit()' to exit)\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
