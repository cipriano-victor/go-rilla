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

	fmt.Printf("Bienvenido a Go-Rilla %s! \n", user.Username)
	fmt.Printf("Presiona Ctrl+C o escribe 'exit' para salir \n")
	repl.Start(os.Stdin, os.Stdout)
}
