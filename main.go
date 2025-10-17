package main

import (
	"flag"
	"fmt"
	"go-rilla/repl"
	"os"
	"os/user"
	"strings"
)

func main() {
	mode := flag.String("mode", string(repl.ModeParser), "execution mode: parser or scanner")
	flag.Parse()

	selectedMode := repl.ModeParser
	switch strings.ToLower(*mode) {
	case string(repl.ModeParser):
		selectedMode = repl.ModeParser
	case string(repl.ModeScanner):
		selectedMode = repl.ModeScanner
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q; valid values are %q or %q\n", *mode, repl.ModeParser, repl.ModeScanner)
		os.Exit(2)
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Welcome to Go-Rilla, %s! \n", currentUser.Username)
	fmt.Printf("Feel free to write any kind of commands \n")
	fmt.Printf("Press Ctrl+C or type 'exit' to leave \n")
	if selectedMode == repl.ModeScanner {
		repl.StartScanner(os.Stdin, os.Stdout)
		return
	}

	repl.StartParser(os.Stdin, os.Stdout)
}
