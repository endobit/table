// Package main shows some features of the table package.
package main

import (
	"fmt"

	"github.com/endobit/table"
	"github.com/endobit/table/sgr"
	"github.com/endobit/table/sgr/color"
)

type rank int

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
	_ = t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-0", Rack: "0", Rank: 0})
	_ = t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-1", Rack: "0", Rank: 1})
	_ = t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-2", Rack: "0", Rank: 2})
	_ = t.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-3", Rack: "0", Rank: 3})
	_ = t.Flush()

	fmt.Println()
	_ = t.FlushJSON()

	fmt.Println()
	_ = t.FlushYAML()

}
