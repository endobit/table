package table

import (
	"bytes"
	"fmt"
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
	t.Flush()

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
	t.Flush()

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
	t.Flush()

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
	t.Flush()

	fmt.Print(buf.String())
	// Output:
	// NAME       PRIORITY
	// Fix bug    8
	// Write docs 3
}
