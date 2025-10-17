package repl

import (
	"bufio"
	"fmt"
	"go-rilla/diag"
	"go-rilla/lexer"
	"go-rilla/parser"
	"go-rilla/token"
	"io"
	"strings"
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
			printDiagnosticsPlain("<repl>", line, ds, out)
		}
		return
	}

	io.WriteString(out, program.String())
	io.WriteString(out, "\n")

	// Util para warnings aún con el parseo correcto
	if ds := l.Diagnostics(); len(ds) > 0 {
		printDiagnosticsPlain("<repl>", line, ds, out)
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
	diags := l.Diagnostics()
	if len(diags) > 0 {
		printDiagnosticsPlain("<repl>", line, diags, out)
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

func printDiagnosticsPlain(filename, src string, diags []diag.Diagnostic, out io.Writer) {
	lines := strings.Split(src, "\n")
	for _, d := range diags {
		fmt.Fprintf(out, "%s:%d:%d: %s %s: %s\n",
			filename, d.Range.Start.Line, d.Range.Start.Column,
			strings.ToLower(d.Level.String()), d.Code, d.Message,
		)

		ln := d.Range.Start.Line
		if ln >= 1 && ln <= len(lines) {
			code := lines[ln-1]
			fmt.Fprintln(out, code)

			runes := []rune(code)
			startCol := clamp(d.Range.Start.Column, 1, len(runes)+1)
			endCol := clamp(d.Range.End.Column, startCol, len(runes)+1)
			length := endCol - startCol
			if length < 1 {
				length = 1
			}

			indent := strings.Repeat(" ", startCol-1)
			underline := indent + "^" + strings.Repeat("~", length-1)
			fmt.Fprintln(out, underline)
		}
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
