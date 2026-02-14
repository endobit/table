// Package table implements a row-based data set that can be printed with a
// header and space aligned columns. The text output looks like a spreadsheet.
package table

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"

	"endobit.io/table/sgr"
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
	Header     []sgr.Param
	EvenRow    []sgr.Param
	OddRow     []sgr.Param
	Empty      []sgr.Param
	Repeat     []sgr.Param
	Annotation []sgr.Param
}

// Table holds a slice of structs that can be Flush()ed as a Text table, or
// encoded as JSON or YAML.
type Table struct {
	rows         []any
	annotations  []annotation
	colors       Colors
	noColor      bool
	writer       io.Writer
	style        style
	fieldToLabel func(string) string
}

type annotation struct {
	index int
	text  string
}

// WithColor is an option setting function for New. It replaces the default set
// of Colors with c.
func WithColor(c *Colors) func(*Table) {
	return func(t *Table) {
		if c != nil {
			t.colors = *c
		}
	}
}

// WithWriter is an option setting function for New. It replaces the default
// io.Writer with w. The io.Writer is used for all Table output.
func WithWriter(w io.Writer) func(*Table) {
	return func(t *Table) {
		t.writer = w
	}
}

// WithLabelFunction is an option setting function for New. This function
// convert struct field names into text header labels. The default behavior is
// to convert the CamelCase field names into UPPER_CASE labels. The "table"
// struct tags can be used to override this.
func WithLabelFunction(fn func(string) string) func(*Table) {
	return func(t *Table) {
		t.fieldToLabel = fn
	}
}

// New returns a new Table. The default settings can be overridden using the
// With* options setting functions. For example: WithColors() can be used to
// replace the default coloring scheme.
func New(opts ...func(*Table)) *Table {
	t := Table{
		writer: os.Stdout,
		colors: Colors{
			Header:     []sgr.Param{sgr.Underline, sgr.Bold},
			Empty:      []sgr.Param{sgr.Faint},
			Repeat:     []sgr.Param{sgr.Faint},
			Annotation: []sgr.Param{sgr.Italic},
		},
		fieldToLabel: camelToUpperSnake,
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

// Write appends the struct a to the table as a row. The current table will be
// flushed if a new struct type is written. If a is not a struct, an error table
// will be added to the output.
func (t *Table) Write(a any) {
	if reflect.TypeOf(a).Kind() != reflect.Struct {
		msg := struct {
			Error error
			Type  string
			Value string
		}{
			Error: ErrNotStruct,
			Type:  fmt.Sprintf("%T", a),
			Value: valueAsString(reflect.ValueOf(a)),
		}

		t.rows = append(t.rows, msg)

		t.Annotate(fmt.Sprintf("skipping non-struct type: %T", a))

		return
	}

	t.rows = append(t.rows, a)
}

// Annotate inserts a string into the table as a row. This is useful for
// inserting comments or other information that is not a struct. The string will
// be printed as-is, without any formatting or coloring.
//
// Annotations are only used in text output, and are ignored in JSON or YAML
// formats.
func (t *Table) Annotate(s string) {
	t.annotations = append(t.annotations, annotation{
		index: len(t.rows),
		text:  s,
	})
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

func maxStringLength(list []string) int {
	maxLen := 0
	for _, s := range list {
		if l := utf8.RuneCountInString(s); l > maxLen {
			maxLen = l
		}
	}

	return maxLen
}

func valueAsString(v reflect.Value) string {
	if v.IsValid() && v.CanInterface() {
		return fmt.Sprintf("%v", v.Interface())
	}

	return ""
}

// camelToUpperSnake converts a CamelCase string to UPPERCASE_SNAKE_CASE,
// supporting Unicode letters.
func camelToUpperSnake(s string) string {
	var (
		b    strings.Builder
		prev rune
	)

	for i, r := range s {
		if i > 0 {
			// Insert underscore if:
			// 1. transition from lower to upper (e.g., "camelCase")
			// 2. transition from letter followed by upper+lower (e.g., "URLValue" -> "URL_Value")
			if unicode.IsLower(prev) && unicode.IsUpper(r) ||
				unicode.IsUpper(prev) && unicode.IsUpper(r) && i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
				b.WriteRune('_')
			}
		}

		b.WriteRune(unicode.ToUpper(r))
		prev = r
	}

	return b.String()
}
