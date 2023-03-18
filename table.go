// Package table implements a row-based data set that can the printed with a
// header and space aligned columns. The text output looks like a spreadsheet.
package table

import (
	"errors"
	"io"
	"os"
	"reflect"

	"golang.org/x/term"

	"github.com/endobit/table.git/sgr"
)

// ErrNotStruct is returned from a Table's Write method if the argument is not a
// struct.
var ErrNotStruct = errors.New("not a struct")

type style int

const (
	textOutput style = iota
	jsonOutput
	yamlOutput
)

// Colors is the set styles/colors to be applied to Table elements.
type Colors struct {
	Header  []sgr.Param
	EvenRow []sgr.Param
	OddRow  []sgr.Param
	Empty   []sgr.Param
	Repeat  []sgr.Param
}

// Table holds a slice of structs that can be Flush()ed as a Text table, or
// encoded as JSON or YAML.
type Table struct {
	rows    []any
	colors  Colors
	noColor bool
	writer  io.Writer
	style   style
}

// WithColor is an option setting function for New. It replaces the default set
// of Colors with c.
func WithColor(c Colors) func(*Table) {
	return func(t *Table) {
		t.colors = c
	}
}

// WithWriter is an option setting function for New. It replaces the default
// io.Writer with w. The io.Writer is used for all Table output.
func WithWriter(w io.Writer) func(*Table) {
	return func(t *Table) {
		t.writer = w
	}
}

// New returns a new Table. The default settings can be overridden using the
// With* options setting functions. For example: WithColors() can be used to
// replace the default coloring scheme.
func New(opts ...func(*Table)) *Table {
	t := Table{
		writer: os.Stdout,
		colors: Colors{
			Header: []sgr.Param{sgr.Underline, sgr.Bold},
			Empty:  []sgr.Param{sgr.Faint},
			Repeat: []sgr.Param{sgr.Faint},
		},
	}

	for _, o := range opts {
		o(&t)
	}

	if !isTerminal(t.writer) {
		t.colors = Colors{}
		t.noColor = true // turn off user Wrapper types
	}

	return &t
}

// Write appends the struct a to the table as a row. It is not an error to Write
// different struct types to the same table.
func (t *Table) Write(a any) error {
	if reflect.TypeOf(a).Kind() != reflect.Struct {
		return ErrNotStruct
	}

	t.rows = append(t.rows, a)

	return nil
}

// Flush writes the table to its writer in its default style.
func (t *Table) Flush() error {
	switch t.style {
	case jsonOutput:
		return t.FlushJSON()
	case yamlOutput:
		return t.FlushYAML()
	default:
		t.FlushText()
	}

	return nil
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(f.Fd()))
}
