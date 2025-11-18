# Go-Rilla: Intérprete de Monkey en Go

## Descripción
Go-Rilla es un intérprete del lenguaje de programación Monkey, implementado en Go. Monkey es un lenguaje de programación simple y dinámico diseñado para enseñar conceptos de interpretación y compilación.

### Disclaimer: 
Este proyecto es una implementación educativa y no está destinado para uso en producción. Está basado en el libro "Writing An Interpreter In Go" de Thorsten Ball.

## Características
- Soporte para tipos de datos básicos: enteros, booleanos, cadenas, arrays y mapas hash.
- Funciones de primera clase, de orden superior y closures.
- Manejo de errores en tiempo de ejecución.
- Evaluación de expresiones aritméticas y lógicas.
- Estructuras de control como condicionales.
- Soporte para declaraciones de variables y asignaciones.
- Funciones Built-In tales como `len`, `first`, `last`, `rest`, `push` y `print`.

## Mejoras
Respecto a la implementación original del libro, Go-Rilla incluye las siguientes mejoras:
- Tracking de líneas y columnas para mensajes de error más precisos en etapas de escaneo y parseo.
- Soporte para numeros de punto flotante y de caracteres Unicode/UTF-8.
- Soporte adicional de operadores de dos caracteres adicionales: `**`, `+=`, `-=`, `&&` y `||`, incluyendo su respectiva precedencia y asociatividad.
- Implementación de operadores postfix `++` y `--`.
- Soporte para errores adicionales en el uso de cadenas tales como escapes inválidos y falta de cierre.
- Implementación de un REPL (Read-Eval-Print Loop) con división por etapas de la interpretación.
- Soporte para bucles `for` y `while`, junto a las declaraciones `continue` y `break`.

## Mejoras Futuras
- De las mencionadas en el libro:
    - Un pretty printer para el AST.
    - Mutabilidad de arrays y mapas hash.
    - Añadir más funciones Built-In.
    - Cachear las hashkeys para mejorar el rendimiento...

- Adicionales: 
    - Bloques `else if`.
    - Operador ternario
    - Comentarios (simples y multilínea)
    - Warnings.
    - Optimizaciones... realmente mucho más!

## Instalación
1. Asegúrate de tener Go instalado en tu sistema.
2. Clona este repositorio:
```bash
git clone https://github.com/cipriano-victor/go-rilla.git
```
3. Navega al directorio del proyecto:
```bash
cd go-rilla
```
4. Ejecuta el intérprete:
```bash
go run main.go
```

## Uso
Por defecto, el intérprete ejecuta un REPL en modo `evaluator`. Puedes escribir código Monkey directamente en la terminal.

Donde podrás elegir entre las diferentes etapas del intérprete: `scanner`, `parser`, y `evaluator`.
```bash
go run main.go -h # Para ver a detalle los flags disponibles
```

Si ya tienes un archivo `.monkey`, usa el flag `-file` para ejecutarlo de una vez y salir:

```bash
go run main.go -file scripts/hello.monkey
```

Puedes combinar `-file` con `-mode` para inspeccionar la salida del `scanner` o del `parser` sin entrar al REPL. Por ejemplo:

```bash
go run main.go -mode scanner -file scripts/hello.monkey
```
