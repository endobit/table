package table

import "github.com/goccy/go-yaml"

// NewYAML returns a Table with YAML as the default for `Flush`.
func NewYAML(opts ...func(*Table)) *Table {
	t := New(opts...)
	t.style = yamlOutput

	return t
}

// FlushYAML flushes the Table data to its io.Writer as YAML.
func (t *Table) FlushYAML() error {
	e := yaml.NewEncoder(t.writer)

	return e.Encode(t.rows)
}
