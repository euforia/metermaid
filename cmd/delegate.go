package main

import (
	"github.com/euforia/metermaid/node"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
)

type GossipDelegate struct {
	node node.Node
	log  *zap.Logger
}

// LocalState satisfies the gossip.Delegate interface
func (del *GossipDelegate) LocalState(join bool) []byte { return nil }

// MergeRemoteState satisfies the gossip.Delegate interface
func (del *GossipDelegate) MergeRemoteState(data []byte, join bool) {}

// NodeMeta satisfies the gossip.Delegate interface
func (del *GossipDelegate) NodeMeta(overhead int) []byte {
	return del.node.MarshalMeta()
}

// NotifyMsg satisfies the gossip.Delegate interface
func (del *GossipDelegate) NotifyMsg(msg []byte) {}

// NotifyJoin satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyJoin(node *memberlist.Node) {
	del.log.Info("node joined",
		zap.String("name", node.Name),
		zap.String("addr", node.Address()),
	)
}

// NotifyLeave satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyLeave(node *memberlist.Node) {}

// NotifyUpdate satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyUpdate(node *memberlist.Node) {
	del.log.Info("node updated", zap.String("addr", node.Address()))
}

func newNode(in *memberlist.Node) *node.Node {
	n := &node.Node{
		Name:    in.Name,
		Address: in.Address(),
	}
	n.UnmarshalMeta(in.Meta)
	return n
}
