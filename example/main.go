// Package main shows some features of the table package.
package main

import (
	"fmt"

	"endobit.io/table"
	"endobit.io/table/sgr"
	"endobit.io/table/sgr/color"
)

type rank int

// Wrap implements the table.wrapper interface for r.
func (r rank) Wrap() sgr.Wrapped {
	return sgr.Wrap(color.Green, r)
}

type host struct {
	Zone    string `table:"ZONE"`
	Cluster string `table:"CLUSTER"`
	Host    string `table:"HOST"`
	Rack    string `table:"RACK,omitempty"`
	Rank    rank   `table:"RANK"`
	Slot    string `table:"SLOT,omitempty"`
}

func main() {
	t := table.New()
	t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-0", Rack: "0", Rank: 0})
	t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-1", Rack: "0", Rank: 1})
	t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-2", Rack: "0", Rank: 2})
	t.Annotate("inline annotation")
	t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-3", Rack: "0", Rank: 3})
	_ = t.Flush()

	fmt.Println()

	_ = t.FlushJSON()

	fmt.Println()

	_ = t.FlushYAML()
}
