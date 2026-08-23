package table

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// AsJSON is an option setting function for New. It sets JSON as the default
// output format for Flush.
func AsJSON() Option {
	return func(t *Table) {
		t.style = jsonOutput
	}
}

// NewJSON returns a Table with JSON as the default for Flush.
//
// Deprecated: Use New(AsJSON()) instead.
func NewJSON(opts ...Option) *Table {
	opts = append([]Option{AsJSON()}, opts...)

	return New(opts...)
}

// FlushJSON flushes the Table data to its io.Writer as JSON.
func (t *Table) FlushJSON() error {
	return json.MarshalWrite(t.writer, t.rows,
		jsontext.WithIndent("    "))
}
