package storage

import (
	"github.com/euforia/gossip"
	"github.com/euforia/metermaid/node"
)

// Nodes implements a node storage interface
type Nodes interface {
	Iter(func(node.Node) error) error
}

// GossipNodes implements the Nodes interface using the underlying gossip
// pool as the backend
type GossipNodes struct {
	pool *gossip.Pool
}

// NewGossipNodes returns a new instance of GossipNodes
func NewGossipNodes(pool *gossip.Pool) *GossipNodes {
	return &GossipNodes{pool: pool}
}

// Iter satisfies the Nodes interface
func (nodes *GossipNodes) Iter(f func(node.Node) error) error {
	var (
		list = nodes.pool.Members()
		err  error
	)
	for _, item := range list {
		n := node.Node{
			Name:    item.Name,
			Address: item.Address(),
		}
		n.UnmarshalMeta(item.Meta)

		if err = f(n); err != nil {
			return err
		}
	}
	return nil
}
