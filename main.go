package main

import (
	"fmt"
	"go-rilla/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Welcome to Go-Rilla, %s! \n", user.Username)
	fmt.Printf("Feel free to write any kind of commands \n")
	fmt.Printf("Press Ctrl+C or type 'exit' to leave \n")
	repl.Start(os.Stdin, os.Stdout)
}
