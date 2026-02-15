package table

import "encoding/json"

// AsJSON is an option setting function for New. It sets JSON as the default
// output format for Flush.
func AsJSON() func(*Table) {
	return func(t *Table) {
		t.style = jsonOutput
	}
}

// NewJSON returns a Table with JSON as the default for Flush.
//
// Deprecated: Use New(AsJSON()) instead.
func NewJSON(opts ...func(*Table)) *Table {
	opts = append([]func(*Table){AsJSON()}, opts...)

	return New(opts...)
}

// FlushJSON flushes the Table data to its io.Writer as JSON.
func (t *Table) FlushJSON() error {
	e := json.NewEncoder(t.writer)
	e.SetIndent("", "    ")

	return e.Encode(t.rows)
}
