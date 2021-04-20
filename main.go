package main

import (
	"fmt"
	"os"
	"os/user"
	"pandey-abhinav/monkey-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hi %s! Monkey programming language!\n", user.Username)
	fmt.Printf("you can give commands below :-\n")
	repl.Start(os.Stdin, os.Stdout)
}
