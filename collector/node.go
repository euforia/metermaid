package collector

import (
	"context"
	"errors"
	"time"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
)

// NodeCollector implements a node runtime collector
type NodeCollector struct {
	node node.Node
	meta []string
}

// Name satisfies the Collector interface
func (nc *NodeCollector) Name() string {
	return "node"
}

// Init satisfies the Collector interface
func (nc *NodeCollector) Init(config *Config) error {
	conf := config.Config
	if tags, ok := conf["meta"]; ok {
		out, err := ifaceSliceToStringSlice(tags)
		if err != nil {
			return err
		}
		nc.meta = out
	}

	if config.Node != nil {
		nc.node = *config.Node
		return nil
	}

	return errors.New("node config required")
}

// Collect satisfies the Collector interface
func (nc *NodeCollector) Collect(context.Context) ([]RunStats, error) {
	rt := RunStats{
		Resource: ResourceNode,
		Start:    time.Unix(0, int64(nc.node.BootTime)),
		Meta:     types.Meta{"node": nc.node.Name},
	}

	for _, k := range nc.meta {
		if val, ok := nc.node.Meta[k]; ok {
			rt.Meta[k] = val
		}
	}

	return []RunStats{rt}, nil
}

func (nc *NodeCollector) Updates() <-chan RunStats { return nil }

func (nc *NodeCollector) Stop() {}
