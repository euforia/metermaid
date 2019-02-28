package collector

import (
	"context"
	"errors"
	"time"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
)

type NodeCollector struct {
	node node.Node
	meta []string
}

func (nc *NodeCollector) Name() string {
	return "node"
}

func (nc *NodeCollector) Init(conf map[string]interface{}) error {
	if tags, ok := conf["meta"]; ok {
		out, err := ifaceSliceToStringSlice(tags)
		if err != nil {
			return err
		}
		nc.meta = out
	}

	if n, ok := conf["node"]; ok {
		if d, ok := n.(node.Node); ok {
			nc.node = d
			return nil
		}
	}
	return errors.New("node config required")
}

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
