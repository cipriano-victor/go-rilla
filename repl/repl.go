package repl

import (
	"bufio"
	"fmt"
	"go-rilla/ast"
	"go-rilla/evaluator"
	"go-rilla/internal/diagprint"
	"go-rilla/lexer"
	"go-rilla/object"
	"go-rilla/parser"
	"go-rilla/token"
	"io"
)

type Mode string

const (
	ModeEvaluator Mode = "evaluator"
	ModeParser    Mode = "parser"
	ModeScanner   Mode = "scanner"
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

func StartEvaluator(in io.Reader, out io.Writer) { startRepl(ModeEvaluator, in, out) }
func StartParser(in io.Reader, out io.Writer)    { startRepl(ModeParser, in, out) }
func StartScanner(in io.Reader, out io.Writer)   { startRepl(ModeScanner, in, out) }

func startRepl(mode Mode, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var env *object.Environment
	if mode == ModeEvaluator {
		env = object.NewEnvironment()
	}
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

		env = processLine(mode, line, out, env)
	}
}

func processLine(mode Mode, line string, out io.Writer, env *object.Environment) *object.Environment {
	switch mode {
	case ModeScanner:
		runScanner(line, out)
	case ModeParser:
		runParser(line, out)
	case ModeEvaluator:
		runEvaluator(line, out, env)
	default:
		runEvaluator(line, out, env)
	}
	return env
}

func StartProgram(line string, out io.Writer) (*lexer.Lexer, *parser.Parser, *ast.Program) {
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		writeDiagnostics(l, p, line, out)
		return nil, nil, nil
	}

	return l, p, program
}

func runEvaluator(line string, out io.Writer, env *object.Environment) {
	l, p, program := StartProgram(line, out)
	if l == nil || p == nil || program == nil {
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}

	writeDiagnostics(l, p, line, out)
}

func runParser(line string, out io.Writer) {
	l, p, program := StartProgram(line, out)
	if l == nil || p == nil || program == nil {
		return
	}

	io.WriteString(out, program.String())
	io.WriteString(out, "\n")

	// Util para warnings aún con el parseo correcto
	writeDiagnostics(l, p, line, out)
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

func writeDiagnostics(l *lexer.Lexer, p *parser.Parser, source string, out io.Writer) {
	if ds := l.Diagnostics(); len(ds) > 0 {
		io.WriteString(out, diagprint.RenderPlain("<repl>", source, ds))
	}
	if pds := p.Diagnostics(); len(pds) > 0 {
		io.WriteString(out, diagprint.RenderPlain("<repl>", source, pds))
	}
}
