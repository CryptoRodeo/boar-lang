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

	fmt.Printf("Hello %s, feel free to type in commands\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
