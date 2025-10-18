package diagprint

import (
	"fmt"
	"strings"

	"go-rilla/diag"
)

// RenderPlain imprime diagnósticos en formato simple (sin ANSI):
// filename:line:col: level CODE: message\n
// <línea de código>\n
// <espacios>^~~~\n
func RenderPlain(filename, src string, diags []diag.Diagnostic) string {
	var out strings.Builder
	lines := strings.Split(src, "\n")
	for _, d := range diags {
		lvl := strings.ToLower(d.Level.String())
		line := clamp(d.Range.Start.Line, 1, len(lines))
		col := max(1, d.Range.Start.Column)
		fmt.Fprintf(&out, "%s:%d:%d: %s %s: %s\n", filename, line, col, lvl, d.Code, d.Message)
		if line-1 < len(lines) {
			code := lines[line-1]
			out.WriteString(code)
			out.WriteString("\n")
			// subrayado
			runes := []rune(code)
			startCol := clamp(d.Range.Start.Column, 1, len(runes)+1)
			var endCol int
			if d.Range.End.Line == d.Range.Start.Line && d.Range.End.Column > 0 {
				endCol = clamp(d.Range.End.Column, startCol, len(runes)+1)
			} else {
				endCol = len(runes) + 1 // si abarca varias líneas, marcamos hasta el fin de la línea
			}
			length := endCol - startCol
			if length < 1 {
				length = 1
			}
			indent := strings.Repeat(" ", startCol-1)
			out.WriteString(indent)
			out.WriteString("^")
			out.WriteString(strings.Repeat("~", length-1))
			out.WriteString("\n")
		}
	}
	return out.String()
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
