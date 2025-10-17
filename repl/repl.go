package repl

import (
	"bufio"
	"fmt"
	"go-rilla/lexer"
	"go-rilla/parser"
	"io"
)

const PROMPT = ">> "
const GORILLA_FACE = `
            __,__
         .-"     "-.
       .'  _   _    '.
      /   (o) (o)     \
     :     .---.       :
     |    /     \      |
     :    \_^_ _/      ;
      \      |        /
       '.    |     .'
         '-.___.-'
		 
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, GORILLA_FACE)
	io.WriteString(out, "Woops! We ran into some gorilla business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
