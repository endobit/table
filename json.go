package table

import "encoding/json"

// NewJSON returns a Table with JSON as the default for `Flush`.
func NewJSON(opts ...func(*Table)) *Table {
	t := New(opts...)
	t.style = jsonOutput

	return t
}

// FlushJSON flushes the Table data to its io.Writer as JSON.
func (t *Table) FlushJSON() error {
	e := json.NewEncoder(t.writer)
	e.SetIndent("", "    ")

	return e.Encode(t.rows)
}
