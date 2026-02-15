package table

import "github.com/goccy/go-yaml"

// AsYAML is an option setting function for New. It sets YAML as the default
// output format for Flush.
func AsYAML() func(*Table) {
	return func(t *Table) {
		t.style = yamlOutput
	}
}

// NewYAML returns a Table with YAML as the default for Flush.
//
// Deprecated: Use New(AsYAML()) instead.
func NewYAML(opts ...func(*Table)) *Table {
	opts = append([]func(*Table){AsYAML()}, opts...)

	return New(opts...)
}

// FlushYAML flushes the Table data to its io.Writer as YAML.
func (t *Table) FlushYAML() error {
	e := yaml.NewEncoder(t.writer)

	return e.Encode(t.rows)
}
