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
	mode := flag.String("mode", string(repl.ModeEvaluator), "execution mode: evaluator, scanner or parser")
	file := flag.String("file", "", "Monkey source file (.monkey) to execute")
	flag.Parse()

	selectedMode := repl.ModeParser
	switch strings.ToLower(*mode) {
	case string(repl.ModeParser):
		selectedMode = repl.ModeParser
	case string(repl.ModeScanner):
		selectedMode = repl.ModeScanner
	case string(repl.ModeEvaluator):
		selectedMode = repl.ModeEvaluator
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q; valid values are %q, %q or %q\n", *mode, repl.ModeParser, repl.ModeScanner, repl.ModeEvaluator)
		os.Exit(2)
	}

	if *file != "" {
		if err := runFile(selectedMode, *file); err != nil {
			fmt.Fprintf(os.Stderr, "could not read %q: %v\n", *file, err)
			os.Exit(1)
		}
		return
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Welcome to Go-Rilla, %s! \n", currentUser.Username)
	fmt.Printf("Feel free to write any kind of commands \n")
	fmt.Printf("Press Ctrl+C/D or type 'exit' to leave \n")
	switch selectedMode {
	case repl.ModeScanner:
		repl.StartScanner(os.Stdin, os.Stdout)
	case repl.ModeParser:
		repl.StartParser(os.Stdin, os.Stdout)
	case repl.ModeEvaluator:
		repl.StartEvaluator(os.Stdin, os.Stdout)
	default:
		repl.StartEvaluator(os.Stdin, os.Stdout)
	}
}

func runFile(mode repl.Mode, path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	repl.RunScript(mode, path, string(source), os.Stdout)
	return nil
}
