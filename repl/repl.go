package repl

import (
	"bufio"
	"fmt"
	"go-rilla/internal/diagprint"
	"go-rilla/lexer"
	"go-rilla/parser"
	"go-rilla/token"
	"io"
)

type Mode string

const (
	ModeParser  Mode = "parser"
	ModeScanner Mode = "scanner"
)

const PROMPT = ">> "
const GORILLA_FACE = `
            __,__
         .-"     "-.
       .'  _   _    '.
      /   (o) (o)     \
     :     .---.       :
     |    / o o \      |
     :    \_ ^_ /      ;
      \      |        /
       '.    |     .'
         '-.___.-'
		 
`

func Start(in io.Reader, out io.Writer) {
	StartParser(in, out)
}

func StartParser(in io.Reader, out io.Writer) {
	startRepl(ModeParser, in, out)
}

func StartScanner(in io.Reader, out io.Writer) {
	startRepl(ModeScanner, in, out)
}

func startRepl(mode Mode, in io.Reader, out io.Writer) {
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

		processLine(mode, line, out)
	}
}

func processLine(mode Mode, line string, out io.Writer) {
	switch mode {
	case ModeScanner:
		runScanner(line, out)
	default:
		runParser(line, out)
	}
}

func runParser(line string, out io.Writer) {
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		if ds := l.Diagnostics(); len(ds) > 0 {
			io.WriteString(out, diagprint.RenderPlain("<repl>", line, ds))
		}
		if pds := p.Diagnostics(); len(pds) > 0 {
			io.WriteString(out, diagprint.RenderPlain("<repl>", line, pds))
		}
		return
	}
	io.WriteString(out, program.String())
	io.WriteString(out, "\n")

	// Util para warnings aún con el parseo correcto
	if ds := l.Diagnostics(); len(ds) > 0 {
		io.WriteString(out, diagprint.RenderPlain("<repl>", line, ds))
	}
	if pds := p.Diagnostics(); len(pds) > 0 {
		io.WriteString(out, diagprint.RenderPlain("<repl>", line, pds))
	}
}

func runScanner(line string, out io.Writer) {
	l := lexer.New(line)
	for {
		tok := l.NextToken()
		fmt.Fprintf(out, "%s\t%q\n", tok.Type, tok.Literal)
		if tok.Type == token.EOF {
			break
		}
	}

	if ds := l.Diagnostics(); len(ds) > 0 {
		io.WriteString(out, diagprint.RenderPlain("<repl>", line, ds))
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
