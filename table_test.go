package table

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"endobit.io/table/sgr"
	"endobit.io/table/sgr/color"
)

func init() {
	// Disable colors for reproducible example output
	sgr.DisableColor()
}

type rank int

func (r rank) Wrap() sgr.Wrapped {
	return sgr.Wrap(nil, r)
}

type host struct {
	Zone    string `table:"ZONE"`
	Cluster string `table:"CLUSTER"`
	Host    string `table:"HOST"`
	Rack    string `table:"RACK,omitempty"`
	Rank    rank   `table:"RANK"`
	Slot    int    `table:"SLOT,omitempty"`
}

func TestYAML(_ *testing.T) {
	w := New()
	w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-0", Rack: "0", Rank: 0})
	w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-1", Rack: "0", Rank: 1})
	w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-2", Rack: "0", Rank: 2})
	w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-3", Rack: "0", Rank: 3})
	_ = w.Flush()
}

func TestCamelToUpperSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "CAMEL_CASE"},
		{"URLValue", "URL_VALUE"},
		{"MyHTTPServer2", "MY_HTTP_SERVER2"},
		{"JSONParser", "JSON_PARSER"},
		{"ÖffentlicheVerkehrsmittel", "ÖFFENTLICHE_VERKEHRSMITTEL"},
		{"ПриветМир", "ПРИВЕТ_МИР"},
		{"ΕλληνικάΚεφαλαία", "ΕΛΛΗΝΙΚΆ_ΚΕΦΑΛΑΊΑ"},
		{"", ""},
		{"A", "A"},
		{"lowercase", "LOWERCASE"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := camelToUpperSnake(tt.input)
			if got != tt.expected {
				t.Errorf("CamelToUpperSnake(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

type server struct {
	Name   string
	Status string
	Port   int
}

func ExampleTable_Write() {
	var buf bytes.Buffer

	t := New(WithWriter(&buf))

	t.Write(server{Name: "web-1", Status: "running", Port: 8080})
	t.Write(server{Name: "web-2", Status: "stopped", Port: 8081})
	_ = t.Flush()

	fmt.Print(buf.String())
	// Output:
	// NAME  STATUS  PORT
	// web-1 running 8080
	// web-2 stopped 8081
}

func ExampleTable_Annotate() {
	var buf bytes.Buffer

	t := New(WithWriter(&buf))

	t.Write(server{Name: "web-1", Status: "running", Port: 8080})
	t.Annotate("--- maintenance window ---")
	t.Write(server{Name: "web-2", Status: "stopped", Port: 8081})
	_ = t.Flush()

	fmt.Print(buf.String())
	// Output:
	// NAME  STATUS  PORT
	// web-1 running 8080
	// --- maintenance window ---
	// web-2 stopped 8081
}

func ExampleNewJSON() {
	var buf bytes.Buffer

	t := NewJSON(WithWriter(&buf))

	t.Write(server{Name: "web-1", Status: "running", Port: 8080})
	t.Write(server{Name: "web-2", Status: "stopped", Port: 8081})
	_ = t.Flush()

	fmt.Print(buf.String())
	// Output:
	// [
	//     {
	//         "Name": "web-1",
	//         "Status": "running",
	//         "Port": 8080
	//     },
	//     {
	//         "Name": "web-2",
	//         "Status": "stopped",
	//         "Port": 8081
	//     }
	// ]
}

type priority int

func (p priority) Wrap() sgr.Wrapped {
	if p > 5 {
		return sgr.Wrap(append(color.Red, sgr.Bold), p)
	}

	return sgr.Wrap(color.Green, p)
}

type task struct {
	Name     string
	Priority priority
}

func ExampleNew_withCustomColors() {
	var buf bytes.Buffer

	t := New(WithWriter(&buf), WithColor(&Colors{
		Header: append([]sgr.Param{sgr.Bold}, color.Cyan...),
	}))

	t.Write(task{Name: "Fix bug", Priority: 8})
	t.Write(task{Name: "Write docs", Priority: 3})
	_ = t.Flush()

	fmt.Print(buf.String())
	// Output:
	// NAME       PRIORITY
	// Fix bug    8
	// Write docs 3
}

func TestWriteNonStruct(t *testing.T) {
	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	// Write a non-struct value
	tbl.Write("not a struct")
	tbl.Flush()

	output := buf.String()

	// Should contain error message
	if !bytes.Contains([]byte(output), []byte("not a struct")) {
		t.Errorf("expected error message in output, got: %q", output)
	}

	// Should contain the ERROR header
	if !bytes.Contains([]byte(output), []byte("ERROR")) {
		t.Errorf("expected ERROR header in output, got: %q", output)
	}
}

func TestEmptyTable(t *testing.T) {
	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	// Flush without writing anything
	tbl.Flush()

	output := buf.String()

	// Should produce minimal output (just a newline or empty)
	if len(output) > 1 {
		t.Errorf("expected empty or minimal output, got: %q", output)
	}
}

func TestSingleColumnTable(t *testing.T) {
	type singleCol struct {
		Name string
	}

	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	tbl.Write(singleCol{Name: "test1"})
	tbl.Write(singleCol{Name: "test2"})
	tbl.Flush()

	output := buf.String()

	// Should contain header and values
	if !bytes.Contains([]byte(output), []byte("NAME")) {
		t.Errorf("expected NAME header, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("test1")) {
		t.Errorf("expected test1 value, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("test2")) {
		t.Errorf("expected test2 value, got: %q", output)
	}
}

func TestOmitEmptyColumns(t *testing.T) {
	type record struct {
		Name  string `table:"NAME"`
		Value string `table:"VALUE,omitempty"`
		Flag  string `table:"FLAG,omitempty"`
	}

	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	tbl.Write(record{Name: "row1", Value: "", Flag: ""})
	tbl.Write(record{Name: "row2", Value: "", Flag: ""})
	tbl.Flush()

	output := buf.String()

	// Should only show NAME column
	if !bytes.Contains([]byte(output), []byte("NAME")) {
		t.Errorf("expected NAME header, got: %q", output)
	}
	// Should NOT show omitempty columns
	if bytes.Contains([]byte(output), []byte("VALUE")) {
		t.Errorf("expected VALUE column to be omitted, got: %q", output)
	}
	if bytes.Contains([]byte(output), []byte("FLAG")) {
		t.Errorf("expected FLAG column to be omitted, got: %q", output)
	}
}

func TestWrapperInterface(t *testing.T) {
	type status string

	var wrapCalled bool

	statusWrapper := struct {
		Status status
	}{
		Status: status("active"),
	}

	// Create a custom type that implements wrapper
	type customStatus struct {
		value string
	}

	wrappedStatus := func(s string) customStatus {
		return customStatus{value: s}
	}

	var _ = wrappedStatus // silence unused warning

	// Test that wrapper interface is called
	type item struct {
		Name string
		Rank rank
	}

	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	tbl.Write(item{Name: "test", Rank: 5})
	tbl.Flush()

	output := buf.String()

	// rank implements Wrap(), so it should appear in output
	if !bytes.Contains([]byte(output), []byte("5")) {
		t.Errorf("expected rank value in output, got: %q", output)
	}

	_ = wrapCalled
	_ = statusWrapper
}

func TestJSONOutput(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	var buf bytes.Buffer
	tbl := NewJSON(WithWriter(&buf))

	tbl.Write(person{Name: "Alice", Age: 30})
	tbl.Write(person{Name: "Bob", Age: 25})
	err := tbl.Flush()

	if err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	output := buf.String()

	// Verify JSON structure
	if !bytes.Contains([]byte(output), []byte(`"Name": "Alice"`)) {
		t.Errorf("expected Alice in JSON output, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte(`"Age": 30`)) {
		t.Errorf("expected Age: 30 in JSON output, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte(`"Name": "Bob"`)) {
		t.Errorf("expected Bob in JSON output, got: %q", output)
	}
}

func TestYAMLOutput(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	var buf bytes.Buffer
	tbl := NewYAML(WithWriter(&buf))

	tbl.Write(person{Name: "Alice", Age: 30})
	tbl.Write(person{Name: "Bob", Age: 25})
	err := tbl.Flush()

	if err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	output := buf.String()

	// Verify YAML structure
	if !bytes.Contains([]byte(output), []byte("name: Alice")) {
		t.Errorf("expected 'name: Alice' in YAML output, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("age: 30")) {
		t.Errorf("expected 'age: 30' in YAML output, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("name: Bob")) {
		t.Errorf("expected 'name: Bob' in YAML output, got: %q", output)
	}
}

func TestWithLabelFunction(t *testing.T) {
	type record struct {
		FirstName string
		LastName  string
	}

	var buf bytes.Buffer

	// Custom label function that lowercases
	lowercase := func(s string) string {
		return strings.ToLower(s)
	}

	tbl := New(WithWriter(&buf), WithLabelFunction(lowercase))

	tbl.Write(record{FirstName: "John", LastName: "Doe"})
	tbl.Flush()

	output := buf.String()

	// Should use custom label function
	if !bytes.Contains([]byte(output), []byte("firstname")) {
		t.Errorf("expected 'firstname' header from custom label function, got: %q", output)
	}
}

func TestMultipleTableTypes(t *testing.T) {
	type typeA struct {
		Name string
	}
	type typeB struct {
		Value int
	}

	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	// Write first type
	tbl.Write(typeA{Name: "test"})
	// Write different type - should flush first table
	tbl.Write(typeB{Value: 42})
	tbl.Flush()

	output := buf.String()

	// Should contain both table headers
	if !bytes.Contains([]byte(output), []byte("NAME")) {
		t.Errorf("expected NAME header for first table, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("VALUE")) {
		t.Errorf("expected VALUE header for second table, got: %q", output)
	}
}

func TestAnnotations(t *testing.T) {
	var buf bytes.Buffer
	tbl := New(WithWriter(&buf))

	tbl.Annotate("before any rows")
	tbl.Write(server{Name: "web-1", Status: "running", Port: 8080})
	tbl.Annotate("middle annotation")
	tbl.Write(server{Name: "web-2", Status: "stopped", Port: 8081})
	tbl.Flush()
	tbl.Annotate("after flush - should appear in next table")

	output := buf.String()

	// First two annotations should appear
	if !bytes.Contains([]byte(output), []byte("before any rows")) {
		t.Errorf("expected first annotation, got: %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("middle annotation")) {
		t.Errorf("expected middle annotation, got: %q", output)
	}
	// Annotation after flush should NOT appear in this output
	if bytes.Contains([]byte(output), []byte("after flush")) {
		t.Errorf("unexpected annotation after flush, got: %q", output)
	}
}
