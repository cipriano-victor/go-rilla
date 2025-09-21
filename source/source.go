package source

// Position representa una ubicación en el archivo fuente.
// Line y Column son 1-based (la primera línea/columna es 1).
// Offset es un desplazamiento en bytes desde el inicio del archivo (0-based).
type Position struct {
	Offset int
	Line   int
	Column int
}

// Range representa un intervalo semiabierto [Start, End) en bytes.
// End es exclusivo. Las líneas/columnas de End pueden calcularse
// a partir de Offset cuando se formatee el diagnóstico.
type Range struct {
	Start Position
	End   Position
}
