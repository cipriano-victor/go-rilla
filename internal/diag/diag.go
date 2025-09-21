package diag

import "go-rilla/source"

// Level indica la severidad del diagnóstico.
type Level int

const (
	Error Level = iota
	Warning
	Note
)

func (l Level) String() string {
	switch l {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Note:
		return "note"
	default:
		return "?"
	}
}

// Diagnostic es una unidad de reporte (mensaje con rango y metadatos).
type Diagnostic struct {
	Level   Level
	Code    string       // p. ej., LEX001, PAR002
	Message string       // mensaje principal
	Hint    string       // sugerencia opcional
	Range   source.Range // [start,end)
}
