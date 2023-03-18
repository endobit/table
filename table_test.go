package table_test

import (
	"testing"

	"github.com/endobit/table.git"
	"github.com/endobit/table.git/sgr"
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

func TestYAML(t *testing.T) {
	w := table.New()
	_ = w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-0", Rack: "0", Rank: 0})
	_ = w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-1", Rack: "0", Rank: 1})
	_ = w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-2", Rack: "0", Rank: 2})
	_ = w.Write(host{Zone: "east", Cluster: "prod", Host: "compute-0-3", Rack: "0", Rank: 3})
	_ = w.Flush()
}
