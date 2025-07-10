package table

import (
	"testing"

	"endobit.io/table/sgr"
)

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
